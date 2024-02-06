package operator

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash"
)

func leafSum(api frontend.API, h hash.FieldHasher, data frontend.Variable) frontend.Variable {
	h.Reset()
	h.Write(data)
	res := h.Sum()

	return res
}

func nodeSum(api frontend.API, h hash.FieldHasher, a, b frontend.Variable) frontend.Variable {
	h.Reset()
	h.Write(a, b)
	res := h.Sum()

	return res
}

func ComputeRoot(api frontend.API, h hash.FieldHasher, leaves []frontend.Variable) frontend.Variable {

	if len(leaves) == 1 {
		return leaves[0]
	}

	if len(leaves)%2 != 0 {
		leaves = append(leaves, leaves[len(leaves)-1])
	}

	var parentNodes []frontend.Variable
	for i := 0; i < len(leaves); i += 2 {
		left := leafSum(api, h, leaves[i])
		right := leafSum(api, h, leaves[i+1])

		nodeHash := nodeSum(api, h, left, right)
		parentNodes = append(parentNodes, nodeHash)
	}

	return ComputeRoot(api, h, parentNodes)
}

func ComputeRootFromPath(api frontend.API, mp *merkle.MerkleProof, h hash.FieldHasher, leaf frontend.Variable) frontend.Variable {

	depth := len(mp.Path) - 1
	sum := leafSum(api, h, mp.Path[0])

	// The binary decomposition is the bitwise negation of the order of hashes ->
	// If the path in the plain go code is 					0 1 1 0 1 0
	// The binary decomposition of the leaf index will be 	1 0 0 1 0 1 (little endian)
	binLeaf := api.ToBinary(leaf, depth)

	for i := 1; i < len(mp.Path); i++ { // the size of the loop is fixed -> one circuit per size
		d1 := api.Select(binLeaf[i-1], mp.Path[i], sum)
		d2 := api.Select(binLeaf[i-1], sum, mp.Path[i])
		sum = nodeSum(api, h, d1, d2)
	}

	return sum
}

func MerkleProofToConstraints(root []byte, proofSet [][]byte) merkle.MerkleProof {
	var path []frontend.Variable
	for j := 0; j < len(proofSet); j++ {
		path = append(path, proofSet[j])
	}

	return merkle.MerkleProof{
		RootHash: root,
		Path:     path,
	}
}
