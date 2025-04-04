package main

// vertex shader (transform vertex and prepare data for fragment shader)
func vShader(v *V4, vn *V4) {
	v.m = 1 //make sure that this is always 1
	//apply combined rotation + translation (although translation will not affect vectors in vns)
	*v = mvMult(vertexRTM, *v)
	*vn = mvMult(vertexRTM, *vn)
	// 3d -> 2d perspective projection (divide x,y by z)
	v.x = v.x / v.z
	v.y = v.y / v.z
	//-> screen space
	*v = mvMult(*T0M, *v)
	lpos := V4{3, 0, -1, 0} //light position
	it := 0xff * dot(*vn, lpos)
	if it < 0 {
		it = 0
	}
	v.m = it //for now, store the light intensity for the vertex in the "magic" coordinate m
}

// fragment shader
func fshader(i, j int, bcs V4, f Face, v0, v1, v2 V4, tex []uint32, pix []byte) {
	if bcs.x < 0 || bcs.y < 0 || bcs.z < 0 {
		return
	}
	if i+j*width >= len(zbuff) || i+j*width < 0 {
		return
	}
	z := bcs.x*v0.z + bcs.y*v1.z + bcs.z*v2.z
	if zbuff[i+j*width] < z {
		return
	}
	zbuff[i+j*width] = z
	//fmt.Printf("z=%v\n", z)
	// fmt.Printf("zval=%v\n", z)
	chv := byte(min(0xff, 0xff-z*0xff/2.5)) // channel val
	if lightingEnabled {
		//bcs = mvMult(lInterpM, bcs)
		it := bcs.x*f.lvals[0] + bcs.y*f.lvals[1] + bcs.z*f.lvals[2] //get lighting at point i,j from barycentric coords
		//if it < 0 {
		//	it = 0
		//}
		chv = byte(it / 4)
	}
	if textureEnabled {
		txidx := bcs.x*f.txidxs[0] + bcs.y*f.txidxs[1] + bcs.z*f.txidxs[2]
		tyidx := bcs.x*f.tyidxs[0] + bcs.y*f.tyidxs[1] + bcs.z*f.tyidxs[2]
		//tidx := int(float64(tstride)*txidx + tyidx*float64(tstride)*float64(tstride))
		xidxint := int(float64(tstride) * txidx)
		yidxint := int(float64(tstride) * tyidx)
		tidx := tstride*xidxint + tstride - yidxint
		//fmt.Printf("texture coords (%v, %v)\n", txidx, tyidx)
		tcol := tex[tidx]
		tcb := bgraToBytes(tcol)
		tb := tcb[0]
		tg := tcb[1]
		tr := tcb[2]
		putpixel(int(i), int(j), byte(float64(chv)*float64(tr)/255), byte(float64(chv)*float64(tg)/255), byte(float64(chv)*float64(tb)/255), 0, pix)
		return
	}
	putpixel(int(i), int(j), 0, chv, 0, 0, pix)
}

func altfShader(i, j int, v V4, tex []uint32, size float64, pix []byte) {
	dx := float64(i) - v.x
	dy := float64(j) - v.y
	z := size - dx*dx - dy*dy
	if z < -100 {
		//return
	}
	if zbuff[i+j*width] < z {
		return
	}
	zbuff[i+j*width] = z
	putpixel(i, j, byte(z), byte(z), byte(z), 0, pix)
	//norm := V4{dx, dy, dz}
}
