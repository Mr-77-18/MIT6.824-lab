package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.
//follow is the example about myself
type Liu_args struct {
	Name string
}

type Liu_reply struct {
	Dial string
}

//the type of version that just have a worker
type Filename struct{
	Files []string
}

type Work_name struct{
	Name string
}


//the type of version that could have more worker
type Args struct {
	Text string
}

type Reply struct{
	Ta Task
}

//task that publich by master 
type Task struct{
	//show what this task is:map or reduce
	//map : file name
	//reduce : reduce number
	Map_reduce int//1 mean map task , mean reduce task
	//map task's source
	Filenames []string

	//reduce task's source 
	N_reduce int //this field's range is 0 -- nreduce-1

	//nreduce come from user
	Nreduce int
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
