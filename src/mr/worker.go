package mr

import "fmt"
import "log"
import "net/rpc"
import "hash/fnv"
import "io/ioutil"
import "sort"
import "os"
import "strconv"
import "encoding/json"

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}


type Intermediate_list struct{
	Intermediate []KeyValue
}

//type ByKey []KeyValue
//
//func (a ByKey) Len() int {return len(a)}
//func (a ByKey) Swap(i , j int) {a[i] , a[j] = a[j] , a[i]}
//func (a ByKey) Less(i , j int) bool {return a[i].Key < a[j].Key}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string){

	// Yor worker implementation here.
	task := Call_requst_task().Ta
	if task.Map_reduce < 0 {// indicate all of the tasks have been done
		return 
	}

	if task.Map_reduce == 1{//it mean it is the map tasks
		//and only the Filenames filed belongs to map tasks
		files := task.Filenames

		nreduce := task.Nreduce

		intermediate_list := [nreduce]Intermediate_list{} 
		intermediate := []KeyValue{}

		var n_reduce int

		for _, filename := range files{
			file , err := os.Open(filename)

			if err != nil{
				log.Fatalf("cannot open %v" , filename)
			}

			content , err := ioutil.ReadAll(file)

			if err != nil{
				log.Fatalf("cannot read %v" , filename)
			}

			file.Close()
			kva := mapf(filename , string(content))

			//n_reduce = ihash(kva.Key)

			intermediate = append(intermediate , kva...)
		}

		//save the key-value to intermediate_list
		sort.Sort(ByKey(intermediate))	
		i := 0
		for i < len(intermediate){
			j := i + 1
			for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key{
				j++
			}
			values := []string{}
			n_reduce = ihash(intermediate[i].Key)
			for k := i ; k < j ; k++{
				intermediate_list = append(intermediate_list[n_reduce] , intermediate[k])
			}

			i = j
		}

		//save the key-value to different file
		i = 0
		for _, interm := range intermediate_list{
			if len(interm) != 0 {
				filename = "mr-" + strconv.Intoa(i)
				file , err = os.Open(filename)
				
				if err != nil{
					log.Fatalf("cannot open %v" , filename)
				}
				//save the key-value in filename,which named is mr-<n-reduce>
				enc := json.NewEncoder(file)
				for _, kv := range interm{
					enc.Encode(&kv)
				}
			}
		}

	}else{//it mean it is the reduce tasks
		//and only the n_reduce filed belongs to reduece tasks
		n_reduce := task.N_reduce	

		//after get the n_reduce , we know open which file to do reduce task,because the file name is mr-<n_reduce>
		filename := "mr-" + strconv.Intoa(n_reduce)

		file , err = os.Open(filename)
		
		if err != nil{
			log.Fatalf("cannot open file %v" , filename)
		}

		//readfile and store them in intermediate(key-value) struct
		intermediate := []KeyValue{}
		dec := json.NewDecoder(file)
		for{
			var kv KeyValue
			if err := dec.Decode(&kv) ; err != nil{
				break
			}
			intermediate = append(intermediate , kv)
		}

		sort.Sort(ByKey(intermediate))

		oname := "mr-out-" + strconv.Intoa(n_reduce)
		ofile , _ := os.Create(oname)

		i := 0
		for i < len(intermediate){
			j := i + 1
			for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key{
				j++
			}
			values := []string{}
			for k := i ; k < j ; k++{
				values = append(values , intermediate[k].Value)
			}
			output := reducef(intermediate[i].Key , values)

			fmt.Fprint(ofile , "%v %v\n" , intermediate[i].Key , output)

			i = j
		}
		ofile.Close()
	}

	// uncomment to send the Example RPC to the coordinator.
	//CallExample()
	
}

func Call_requst_task() Reply {
	args := Args{}
	args.Text = "i request a task for map or reduce"

	reply := Reply{}

	ok :=  call("Coordinator.Request_task" , &args , &reply)
	if ok{
	}else{
		fmt.Printf("call request task fail\n")
	}

	return reply
}

//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}


//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

