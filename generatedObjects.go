package main

import (
	"math"
)

// creates vertices of cube from given centre and half diagonal vectors.
// o is centre and u is half diagonal vector
func cuboid(o, u [][]float64) object {
	cu := object{}
	cu.center = vector(o[0][0], o[1][0], o[2][0], 1)

	// creating cube parallel to axes by default
	cu.vertices = make([][][]float64, 8)
	cu.camtices = make([][][]float64, 8)
	cu.vertices[0] = vector(u[0][0], u[1][0], u[2][0], 0)
	cu.vertices[1] = vector(u[0][0], u[1][0], -u[2][0], 0)
	cu.vertices[2] = vector(-u[0][0], u[1][0], -u[2][0], 0)
	cu.vertices[3] = vector(-u[0][0], u[1][0], u[2][0], 0)
	cu.vertices[4] = vector(u[0][0], -u[1][0], u[2][0], 0)
	cu.vertices[5] = vector(u[0][0], -u[1][0], -u[2][0], 0)
	cu.vertices[6] = vector(-u[0][0], -u[1][0], -u[2][0], 0)
	cu.vertices[7] = vector(-u[0][0], -u[1][0], u[2][0], 0)
	for i := 0; i < 8; i++ {
		cu.vertices[i] = matAdd(cu.vertices[i], o)
		cu.camtices[i] = matAdd(cu.vertices[i], o)
	}
	cu.triangles = make([]triangle, 12)
	for i := 0; i < 12; i++ {
		cu.triangles[i] = triangle{}
	}
	cu.triangles[0].create(&cu.vertices[0], &cu.vertices[1], &cu.vertices[2], &cu.camtices[0], &cu.camtices[1], &cu.camtices[2])
	cu.triangles[1].create(&cu.vertices[0], &cu.vertices[2], &cu.vertices[3], &cu.camtices[0], &cu.camtices[2], &cu.camtices[3])
	cu.triangles[2].create(&cu.vertices[4], &cu.vertices[7], &cu.vertices[6], &cu.camtices[4], &cu.camtices[7], &cu.camtices[6])
	cu.triangles[3].create(&cu.vertices[4], &cu.vertices[6], &cu.vertices[5], &cu.camtices[4], &cu.camtices[6], &cu.camtices[5])
	cu.triangles[4].create(&cu.vertices[0], &cu.vertices[4], &cu.vertices[5], &cu.camtices[0], &cu.camtices[4], &cu.camtices[5])
	cu.triangles[5].create(&cu.vertices[0], &cu.vertices[5], &cu.vertices[1], &cu.camtices[0], &cu.camtices[5], &cu.camtices[1])
	cu.triangles[6].create(&cu.vertices[1], &cu.vertices[5], &cu.vertices[6], &cu.camtices[1], &cu.camtices[5], &cu.camtices[6])
	cu.triangles[7].create(&cu.vertices[1], &cu.vertices[6], &cu.vertices[2], &cu.camtices[1], &cu.camtices[6], &cu.camtices[2])
	cu.triangles[8].create(&cu.vertices[2], &cu.vertices[6], &cu.vertices[7], &cu.camtices[2], &cu.camtices[6], &cu.camtices[7])
	cu.triangles[9].create(&cu.vertices[2], &cu.vertices[7], &cu.vertices[3], &cu.camtices[2], &cu.camtices[7], &cu.camtices[3])
	cu.triangles[10].create(&cu.vertices[0], &cu.vertices[7], &cu.vertices[4], &cu.camtices[0], &cu.camtices[7], &cu.camtices[4])
	cu.triangles[11].create(&cu.vertices[0], &cu.vertices[3], &cu.vertices[7], &cu.camtices[0], &cu.camtices[3], &cu.camtices[7])
	
	return cu
}

// returns sphere object with o as centre, r as radius and made of n*n triangles
func sphere(o [][]float64, r float64, n int) object {
	sp := object{}
	sp.center = vector(o[0][0], o[1][0], o[2][0], 1)

	dtheta := math.Pi / float64(n)
	dphi := dtheta * 2
	var theta, phi float64
	for j := 0; j < n+1; j++ {
		for i := 0; i < n; i++ {
			vertex := vector(r*math.Sin(theta)*math.Cos(phi), r*math.Cos(theta), r*math.Sin(theta)*math.Sin(phi), 0) // 4th col is 0 cuz o has 1 there
			sp.vertices = append(sp.vertices, matAdd(vertex, o))
			sp.camtices = append(sp.camtices, vector(vertex[0][0], vertex[1][0], vertex[2][0], vertex[3][0]))
			phi += dphi
		}
		theta += dtheta
	}

	sp.triangles = make([]triangle, n*(2*n))
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			sp.triangles[j*n*2+i*2] = triangle{}
			sp.triangles[j*n*2+i*2].create(
				&sp.vertices[j*n+i],
				&sp.vertices[j*n+(i+1)%n],
				&sp.vertices[(j+1)*n+i],
				&sp.camtices[j*n+i],
				&sp.camtices[j*n+(i+1)%n],
				&sp.camtices[(j+1)*n+i],
			)
			sp.triangles[j*n*2+i*2+1] = triangle{}
			sp.triangles[j*n*2+i*2+1].create(
				&sp.vertices[(j+1)*n+i],
				&sp.vertices[j*n+(i+1)%n],
				&sp.vertices[(j+1)*n+(i+1)%n],
				&sp.camtices[(j+1)*n+i],
				&sp.camtices[j*n+(i+1)%n],
				&sp.camtices[(j+1)*n+(i+1)%n],
			)
		}
	}
	return sp
}
