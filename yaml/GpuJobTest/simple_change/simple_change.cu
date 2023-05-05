#include <stdio.h>
#include <cuda_runtime.h>
__global__ void add(int *a) {
        a[0]=a[0]*2;
}

int main() {
    int N = 1;
    int *a;
    int *dev_a; 

    a = (int*)malloc(N * sizeof(int)); 
    cudaMalloc((void**)&dev_a, N * sizeof(int));   
    
    a[0]=100;
    cudaMemcpy(dev_a, a, N * sizeof(int), cudaMemcpyHostToDevice);
    printf("a[0] = %d\n", a[0]);
    add<<<1, 1>>>(dev_a);
    cudaMemcpy(a, dev_a, N * sizeof(int), cudaMemcpyDeviceToHost);
    printf("a[0] = %d\n", a[0]);

    free(a);
    cudaFree(dev_a);
    return 0;
}

