#!/bin/bash

#SBATCH --job-name=performance        
#SBATCH --partition=dgx2      
#SBATCH -N 1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=1     
#SBATCH --gres=gpu:1          
#SBATCH --output=result.out
#SBATCH --error=result.err

ulimit -s unlimited
ulimit -l unlimited

module load gcc/8.3.0 cuda/10.1.243-gcc-8.3.0

nvcc -o performance performance.cu
./performance
