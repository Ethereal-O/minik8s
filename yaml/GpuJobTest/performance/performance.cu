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

void judge_result(int *mat_cpu,int *mat_gpu,int rowNum,int colNum){
	for(int i=0;i<rowNum;i++){
                for(int j=0;j<colNum;j++){
                        if(mat_cpu[i*colNum+j]!=mat_gpu[i*colNum+j]){
                            printf("result error!\n");
                            return;
                        }
                }
        }
    printf("result pass!\n");
}



int main(){
	int m=1<<8,n=1<<8,p=1<<8;
	int thread_per_block_x=1<<4;
	int thread_per_block_y=1<<4;
	dim3 block((p+thread_per_block_x-1)/thread_per_block_x,
	            (m+thread_per_block_y-1)/thread_per_block_y);
    dim3 thread(thread_per_block_x,thread_per_block_y);
	int *mat_a,*mat_b,*mat_gpu,*mat_cpu;
	int *dmat_a,*dmat_b,*dmat_c;
	clock_t cpu_start,cpu_end,gpu_start,gpu_end;

	mat_a=(int*)calloc(m*n,sizeof(int));
	mat_b=(int*)calloc(n*p,sizeof(int));
	mat_gpu=(int*)calloc(m*p,sizeof(int));
	mat_cpu=(int*)calloc(m*p,sizeof(int));

	cudaMalloc((void**)&dmat_a,m*n*sizeof(int));
	cudaMalloc((void**)&dmat_b,n*p*sizeof(int));
	cudaMalloc((void**)&dmat_c,m*p*sizeof(int));

	mat_init(mat_a,m,n);
	mat_init(mat_b,n,p);

    cpu_start=clock();
    for(int i=0;i<m;i++){
        for(int j=0;j<p;j++){
            int res=0;
            for(int k=0;k<n;k++){
               res+=mat_a[i*n+k]*mat_b[k*p+j];
            }
            mat_cpu[i*p+j]=res;
        }
    }
    cpu_end=clock();

	cudaMemcpy(dmat_a,mat_a,m*n*sizeof(int),cudaMemcpyHostToDevice);
	cudaMemcpy(dmat_b,mat_b,n*p*sizeof(int),cudaMemcpyHostToDevice);

	gpu_start=clock();
    mat_mul<<<block,thread>>>(dmat_a,dmat_b,dmat_c,m,p,n);

    cudaMemcpy(mat_gpu,dmat_c,m*p*sizeof(int),cudaMemcpyDeviceToHost);
	cudaDeviceSynchronize();
	gpu_end=clock();

	printf("cpu : %d %d<-> gpu: %d %d\n",cpu_end,cpu_start,gpu_end,gpu_start);
	judge_result(mat_cpu,mat_gpu,m,p);


	free(mat_a);
	free(mat_b);
	free(mat_gpu);
	free(mat_cpu);

	cudaFree(dmat_a);
	cudaFree(dmat_b);
	cudaFree(dmat_c);
}