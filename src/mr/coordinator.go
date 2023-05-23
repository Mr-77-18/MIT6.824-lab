//the version has one or more worker
package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"
import "fmt"


//go language , we could see that Coordinator is a c++ class , which have the member values and member function . what is different is that go use function type to define a member function
type Coordinator struct {
	// Your definitions here.
	//indicate file whether is managed,when the file is done , we should remove it in Files
	Files []string

	Task_done int//one indicate a range of Files to publich,such as Task_done=1 , mean 2 files have been done,and "2" could change by user

	N_reduce int//init with n_reduce come from user , and then reduce after pulish

	Nreduce int//come from user,and should not be changed
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) Request_task(args *Args, reply *Reply) error {
	//version that have serval worker , and this funtion publich task to the worker
	
	task_publish := Task{}
	task_publish.Nreduce = c.Nreduce

	if c.Task_done * 2 <= len(c.Files){//it mean the map task have not done yet

		task_publish.Map_reduce = 1
		begin_index = c.Task_done * 2
		end_index = c.Task_done * 2 + 2
		c.Task_done = c.Task_done + 1

		if (end_index) > len(c.Files){
			end_index := len(c.Files)
		}
		task_publish.Filenames = c.Files[begin_index : end_index]


	}else{//it mean the map task had done and should publich reduce task
		task_publish.Map_reduce = 0
		if c.Nreduce == 0 {//all reduce task had been done
			//tell the caller all task have been done
			task_publish.Map_reduce = -1
			return nil
		}else{
			task_publish.N_reduce = c.N_reduce
			c.N_reduce = c.N_reduce - 1
		}
	}

	reply.Ta = task_publish

	return nil
}


// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)//e mean error,l mean listen
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.  
	//when the Nreduce is 0 , indicate all task is done
	if c.N_reduce == 0 {
		ret = true
	}

	return ret
}

//
// create a Coordinator.Coordinator have member values and member function
// how to find the member function , we just need to see what type of the function[the type if local in : func (<type>) function_name(<args>) <return values>
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.
	c.Files = files//remaind the files name in the struct,when the Fils is empty, it indicate that the work is done
	c.Task_done = 0

	c.Nreduce = nReduce
	c.N_reduce = nReduce

	c.server()
	return &c
}
