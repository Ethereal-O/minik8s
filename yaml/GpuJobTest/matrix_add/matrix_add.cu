#include <stdio.h>
#include <iostream>
#include <cuda_runtime.h>
#include <time.h>

__global__ void mat_add(int *dmat_a,int* dmat_b,int *dmat_c,int rowNum,int colNum){
    int idx_x=blockIdx.x*blockDim.x+threadIdx.x;
    int idx_y=blockIdx.y*blockDim.y+threadIdx.y;
	int idx=idx_x+idx_y*colNum;
	if(idx<rowNum*colNum){
	    dmat_c[idx]=dmat_a[idx]+dmat_b[idx];
	}
}

int create_rand(){
        int ret = rand() % 10 + 1;
	return ret;
}

void mat_init(int *mat,int rowNum,int colNum){
    for(int i=0;i<rowNum;i++){
    	for(int j=0;j<colNum;j++){
    		mat[i*colNum+j]=create_rand();
    	}
    }
}

void print_result(int *mat,int m,int n){
	for(int i=0;i<m;i++){
                for(int j=0;j<n;j++){
                        printf("%d\t",mat[i*n+j]);
                }
		printf("\n");
        }
}

int main(){
	int m=1<<3,n=1<<4;
	int thread_per_block_x=1<<3;
	int thread_per_block_y=1<<3;
	dim3 block((n+thread_per_block_x-1)/thread_per_block_x,
	            (m+thread_per_block_y-1)/thread_per_block_y);
    dim3 thread(thread_per_block_x,thread_per_block_y);
	int *mat_a,*mat_b,*mat_c;
	int *dmat_a,*dmat_b,*dmat_c;

	mat_a=(int*)calloc(m*n,sizeof(int));
	mat_b=(int*)calloc(m*n,sizeof(int));
	mat_c=(int*)calloc(m*n,sizeof(int));

	cudaMalloc((void**)&dmat_a,m*n*sizeof(int));
	cudaMalloc((void**)&dmat_b,m*n*sizeof(int));
	cudaMalloc((void**)&dmat_c,m*n*sizeof(int));

	mat_init(mat_a,m,n);
	mat_init(mat_b,m,n);

	cudaMemcpy(dmat_a,mat_a,m*n*sizeof(int),cudaMemcpyHostToDevice);
	cudaMemcpy(dmat_b,mat_b,m*n*sizeof(int),cudaMemcpyHostToDevice);

    mat_add<<<block,thread>>>(dmat_a,dmat_b,dmat_c,m,n);

    cudaMemcpy(mat_c,dmat_c,m*n*sizeof(int),cudaMemcpyDeviceToHost);
	cudaDeviceSynchronize();

    printf("mat_a:\n");
	print_result(mat_a,m,n);
    printf("mat_b:\n");
	print_result(mat_b,m,n);
    printf("mat_c:\n");
	print_result(mat_c,m,n);

	free(mat_a);
	free(mat_b);
	free(mat_c);

	cudaFree(dmat_a);
	cudaFree(dmat_b);
	cudaFree(dmat_c);
}
