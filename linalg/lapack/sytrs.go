
// Copyright (c) Harri Rautila, 2012

// This file is part of go.opt/linalg package. It is free software, distributed
// under the terms of GNU Lesser General Public License Version 3, or any later
// version. See the COPYING tile included in this archive.

package lapack

import (
	"github.com/hrautila/go.opt/linalg"
	"github.com/hrautila/go.opt/matrix"
	"errors"
)

/*
 Solves a real or complex symmetric set of linear equations,
 given the LDL^T factorization computed by sytrf() or sysv().

 Sytrs(A, ipiv, B, uplo=PLower, n=A.Rows, nrhs=B.Cols,
 ldA=max(1,A.Rows), ldB=max(1,B.Rows), offsetA=0, offsetB=0)

 PURPOSE
 Solves
  A*X = B
 
 where A is real or complex symmetric and n by n,
 and B is n by nrhs.  On entry, A and ipiv contain the
 factorization of A as returned by Sytrf() or Sysv().  On exit, B is
 replaced by the solution.

 ARGUMENTS
  A         float or complex matrix
  B         float or complex matrix.  Must have the same type as A.
  ipiv      int vector

 OPTIONS
  uplo      PLower or PUpper
  n         nonnegative integer.  If negative, the default value is used.
  nrhs      nonnegative integer.  If negative, the default value is used.
  ldA       positive integer.  ldA >= max(1,n).  If zero, the default
            value is used.
  ldB       nonnegative integer.  ldB >= max(1,n).  If zero, the
            default value is used.
  offsetA   nonnegative integer
  offsetB   nonnegative integer;

*/
func Sytrs(A, B matrix.Matrix, ipiv []int32, opts ...linalg.Option) error {
	pars, err := linalg.GetParameters(opts...)
	if err != nil {
		return err
	}
	ind := linalg.GetIndexOpts(opts...)
	if ind.N < 0 {
		ind.N = A.Rows()
		if ind.N != A.Cols() {
			return errors.New("A not square")
		}
	}
	if ind.Nrhs < 0 {
		ind.Nrhs = B.Cols()
	}
	if ind.N == 0 || ind.Nrhs == 0 {
		return nil
	}
	if ind.LDa == 0 {
		ind.LDa = max(1, A.Rows())
	}
	if ind.LDa < max(1, ind.N) {
		return errors.New("lda")
	}
	if ind.LDb == 0 {
		ind.LDb = max(1, B.Rows())
	}
	if ind.LDb < max(1, ind.N) {
		return errors.New("ldb")
	}
	if ind.OffsetA < 0 {
		return errors.New("offsetA")
	}
	sizeA := A.NumElements()
	if sizeA < ind.OffsetA+(ind.N-1)*ind.LDa+ind.N {
		return errors.New("sizeA")
	}
	if ind.OffsetB < 0 {
		return errors.New("offsetB")
	}
	sizeB := B.NumElements()
	if sizeB < ind.OffsetB+(ind.Nrhs-1)*ind.LDb+ind.N {
		return errors.New("sizeB")
	}
	if ipiv != nil && len(ipiv) < ind.N {
		return errors.New("size ipiv")
	}
	info := -1
	switch A.(type) {
	case *matrix.FloatMatrix:
		Aa := A.FloatArray()
		Ba := B.FloatArray()
		uplo := linalg.ParamString(pars.Uplo)
		info = dsytrs(uplo, ind.N, ind.Nrhs, Aa[ind.OffsetA:], ind.LDa, ipiv,
			Ba[ind.OffsetB:], ind.LDb)
	case *matrix.ComplexMatrix:
	}
	if info != 0 {
		return errors.New("sytrs failed")
	}
	return nil
}


// Local Variables:
// tab-width: 4
// End:
