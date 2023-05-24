package mr

import "fmt"
import "log"
import "net/rpc"
import "hash/fnv"
import "io/ioutil"
import "sort"
import "os"

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

type ByKey []KeyValue

func (a ByKey) Len() int {return len(a)}
func (a ByKey) Swap(i , j int) {a[i] , a[j] = a[j] , a[i]}
func (a ByKey) Less(i , j int) bool {return a[i].Key < a[j].Key}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}


//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string){

	// Yo
	r worker implementation here.

	// uncomment to send the Example RPC to the coordinator.
	//CallExample()
	File_name := Call_get_filename()//to get the file name in the coordinator
	//do the work
	
	//when the work is done , we need to call the rpc function : Work_done to notic the coordinator that the map and reduce functio is done
	//the plugin has been loaded
	intermediate := []KeyValue{}
	for _, filename := range File_name.Files{
		file , err := os.Open(filename)

		if err != nil{
			log.Fatalf("cannot open %v" , filename)
		}

		content , err := ioutil.ReadAll(file)

		if err != nil{

			log.Fatalf("cannot  read %v" , filename)
		}

		file.Close()

		kva := mapf(filename , string(content))

		intermediate = append(intermediate , kva...)
	}

	sort.Sort(ByKey(intermediate))
	oname := "mr-out-0"
	ofile, _ := os.Create(oname)

	i:=0
	for i < len(intermediate){
		j := i + 1
		for j< len(intermediate) && intermediate[j].Key == intermediate[i].Key{
			j++
		}
		values := []string{}
		for k := i ; k < j ; k++ {
			values = append(values , intermediate[k].Value)
		}

		output := reducef(intermediate[i].Key , values)

		fmt.Fprintf(ofile , "%v %v\n" , intermediate[i].Key , output)

		i = j
	}

	ofile.Close()

	Call_work_done()
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

func Call_get_filename() Filename{

	// declare an argument structure.
	args := Work_name{}

	args.Name = "work1"

	// fill in the argument(s).

	// declare a reply structure.
	reply := Filename{}


	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Get_filename", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("filename read complete\n")
	} else {
		fmt.Printf("read filename error , call failed!\n")
	}

	return reply
}



func Call_work_done() {
	args := ExampleArgs{}
	reply := ExampleReply{}

	ok := call("Coordinator.Work_done" , &args , &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("the work is done\n")
	} else {
		fmt.Printf("read filename error , call failed!\n")
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

