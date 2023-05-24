in lab 1 , we should know that:
1. coordinator.go is called by mrcoordinator.go
2. worker.go is called by mrworker.go

how to run the example is that:
1. run the mrcoordinator.go by pass it the args about the input files
2. run the mrworker.go by pass it the args about the map and reduce function(wc.so)

and what should need to know is that how the worker.go and coordinator works
1. coordinator should split the input files several fragment, which should pass to the worker to do map and reduce
2. worker is do the real task: map and reduce , but this tasks should ask  for coordinator by rpc

note that the worker just need the file name?

**the final achive of the lab is that** 
1. worker: use rpc to request the task(map or reduce)
```shell
task = request_task()
if task == map {
	map()
	sort()//different key-value save in different file
	save_key_val()//filename should be mr-work_number-<reduce_number>,reduce_number come from function ihash(key)/nreduce in file worker.go
}else{
	reduce_number = task.reduce_number
	readfile()//read file name in format mr-*-reduce_number
	sort()
	reduce(
	save_final_output()
}
`````

2. master:publich task(map or reduce) , but reduce task should publiched until maps tasks are done
	1. judgment whether the map tasks is done
	2. publich the task and then reduce the number of the task(more safety version should confirm that the task is complete be done , but this versio not)
	3. when all tasks is done , we need to end the process(it may relate to the Done() function int coordinator.go)
```shell
if map task not done yet {
	publish map task
	map task's number reduce
}else{
	publish reduce task
	reduce  task's number reduce 
	if all reduce tasks done {
		indicate the mapreduce job is done(such as use a mark to indicate that , and judg the job is done by the mark in function Done() in file coordinator.go)
	}
}
`````


