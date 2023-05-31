启动两个Node
1: make master
1: ./kubectl apply -f ./yaml/PodTest/node1.yaml
3: make worker
3: ./kubectl apply -f ./yaml/PodTest/node3.yaml
1: ./kubectl get -t Node
启动一个Pod
1: ./kubectl apply -f ./yaml/PodTest/pod1.yaml
1: ./kubectl get -t Pod
3: docker ps
检查hostPort是否已开放
3: curl localhost:8888
检查Pod是否共享Network Namespace
3: docker exec (curl) curl localhost:80
检查容器的资源占用量
3: docker stats
检查Volume是否可以共享文件
3: docker exec (nginx) touch /usr/share/nginx/html/files/hello.txt
3: docker exec (nginx) echo "hello" > /usr/share/nginx/html/files/hello.txt
3: docker exec (curl) cat /usr/share/nginx/html/files/hello.txt
检查Pod的容错性
3: docker kill (nginx)
3: ./kubectl get -t Pod
3: docker ps
3: curl localhost:8888