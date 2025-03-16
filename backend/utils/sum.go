package utils

import "math"

func calculateTake(price1 float64, price2 float64, listener []int, sender []int) float64 {
    
    x1, x2, y1, y2:= listener[0], sender[0], listener[1], sender[1]
    needX := math.Abs(float64(x2) - float64(x1))
    needY := math.Abs(float64(y2) - float64(y1))

    needLow := math.Sqrt(needX*needX + needY*needY)

    var take float64
    var ost float64

    if price2 <= price1*2 {
       take = math.Floor(needLow/127) * price2
       ost = math.Mod(needLow, 127)
    } else {
       take = math.Floor(needLow/63) * price1 * 2
       ost = math.Mod(needLow, 63)
    }

    if ost <= 63 {
       take += price1
    } else {
       take += price2
    }

    return take
}