package zydis

func Ok(status uint32) bool {
	return int32(status) >= 0
}

//func Failed(status Status) bool {
//	return int32(status) < 0
//}
