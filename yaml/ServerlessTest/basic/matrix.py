def rank(d):
    A = np.array([[1, 2, 3, 4], [2, 4, 6, 8], [3, 6, 9, 12], [4, 8, 12, 16]])
    rank = np.linalg.matrix_rank(A)
    print("the rank of A is ", rank)
    d['rank'] = str(rank)
    return d