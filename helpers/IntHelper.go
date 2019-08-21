package helpers

func IfUint8(condition bool, _then, _else uint8) uint8 {
    if condition {
        return _then
    }
    return _else
}

func MaxUint8(ints ...uint8) uint8 {
    var max uint8;
    for _, v := range ints {
        if v > max {
            max = v
        }
    }
    return max;
}