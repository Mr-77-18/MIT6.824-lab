//the version has just one works
package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"


//go language , we could see that Coordinator is a c++ class , which have the member values and member function . what is different is that go use function type to define a member function
type Coordinator struct {
	// Your definitions here.
	//indicate file whether is managed,when the file is done , we should remove it in Files
	Files []string
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

func (c *Coordinator) Liu_example(args *Liu_args, reply *Liu_reply) error {
	reply.Dial = "hello , nice to meet you" + " " + args.Name
	return nil
}


func (c *Coordinator) Get_filename(args *Work_name , reply *Filename) error{
	reply.Files = c.Files
	return nil
}

func (c *Coordinator) Work_done( args *ExampleArgs, reply *ExampleReply )  error{
	var empty []string
	c.Files = empty
	return nil
}


//
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
	if len(c.Files) == 0 {//when the split is empty , indicates the work is done
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


	c.server()
	return &c
}
