#include <stdio.h>
#include <iostream>
#include <cuda_runtime.h>
#include <time.h>

__global__ void mat_mul(int *dmat_a,int* dmat_b,int *dmat_c,int rowNum,int colNum,int midNum){
    int idx_x=blockIdx.x*blockDim.x+threadIdx.x;
    int idx_y=blockIdx.y*blockDim.y+threadIdx.y;
	int idx=idx_x+idx_y*colNum;
	if(idx<rowNum*colNum){
	    int res=0;
	    for(int i=0;i<midNum;i++){
	        res+=dmat_a[idx_y*midNum+i]*dmat_b[i*colNum+idx_x];
	    }
	    dmat_c[idx]=res;
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

void print_result(int *mat,int rowNum,int colNum){
	for(int i=0;i<rowNum;i++){
                for(int j=0;j<colNum;j++){
                        printf("%d\t",mat[i*colNum+j]);
                }
		printf("\n");
        }
}

int main(){
	int m=1<<5,n=1<<6,p=1<<3;
	int thread_per_block_x=1<<3;
	int thread_per_block_y=1<<3;
	dim3 block((p+thread_per_block_x-1)/thread_per_block_x,
	            (m+thread_per_block_y-1)/thread_per_block_y);
    dim3 thread(thread_per_block_x,thread_per_block_y);
	int *mat_a,*mat_b,*mat_c;
	int *dmat_a,*dmat_b,*dmat_c;

	mat_a=(int*)calloc(m*n,sizeof(int));
	mat_b=(int*)calloc(n*p,sizeof(int));
	mat_c=(int*)calloc(m*p,sizeof(int));

	cudaMalloc((void**)&dmat_a,m*n*sizeof(int));
	cudaMalloc((void**)&dmat_b,n*p*sizeof(int));
	cudaMalloc((void**)&dmat_c,m*p*sizeof(int));

	mat_init(mat_a,m,n);
	mat_init(mat_b,n,p);

	cudaMemcpy(dmat_a,mat_a,m*n*sizeof(int),cudaMemcpyHostToDevice);
	cudaMemcpy(dmat_b,mat_b,n*p*sizeof(int),cudaMemcpyHostToDevice);

    mat_mul<<<block,thread>>>(dmat_a,dmat_b,dmat_c,m,p,n);

    cudaMemcpy(mat_c,dmat_c,m*p*sizeof(int),cudaMemcpyDeviceToHost);
	cudaDeviceSynchronize();

    printf("mat_a:\n");
	print_result(mat_a,m,n);
    printf("mat_b:\n");
	print_result(mat_b,n,p);
    printf("mat_c:\n");
	print_result(mat_c,m,p);

	free(mat_a);
	free(mat_b);
	free(mat_c);

	cudaFree(dmat_a);
	cudaFree(dmat_b);
	cudaFree(dmat_c);
}
