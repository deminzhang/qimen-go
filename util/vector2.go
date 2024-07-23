package util

import (
	"math"
)

type Vec2[T Numeric] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

func (this *Vec2[T]) Equal(v Vec2[T]) bool {
	return *this == v
}

func (this *Vec2[T]) Set(x, y T) {
	this.X = x
	this.Y = y
}

func (this *Vec2[T]) Length() float32 {
	return float32(math.Sqrt(float64((this.X * this.X) + (this.Y * this.Y))))
}

func (this *Vec2[T]) Clone() Vec2[T] {
	return NewVec2[T](this.X, this.Y)
}

func (this *Vec2[T]) Add(v Vec2[T]) {
	this.X += v.X
	this.Y += v.Y
}

func (this *Vec2[T]) Sub(v Vec2[T]) {
	this.X -= v.X
	this.Y -= v.Y
}

func (this *Vec2[T]) Multiply(scalar T) {
	this.X *= scalar
	this.Y *= scalar
}

func (this *Vec2[T]) Divide(scalar T) {
	if scalar == 0 {
		//panic("v/0！")
		this.X = T(math.Inf(1))
		this.Y = T(math.Inf(1))
		return
	}
	this.X /= scalar
	this.Y /= scalar
}

// 向量：分向量乘
func (this *Vec2[T]) Scale(v Vec2[T]) {
	this.X *= v.X
	this.Y *= v.Y
}

// 向量：点积
func (this *Vec2[T]) Dot(v Vec2[T]) T {
	return this.X*v.X + this.Y*v.Y
}

// 向量：长度
func (this *Vec2[T]) Magnitude() T {
	return T(math.Sqrt(float64(this.X*this.X + this.Y*this.Y)))
}

// 向量：长度平方
func (this *Vec2[T]) SqrMagnitude() T {
	return this.X*this.X + this.Y*this.Y
}

// 向量：单位化 (0向量禁用)
func (this *Vec2[T]) Normalize() {
	l := this.Magnitude()
	if l == 0 {
		this.Set(1, 0)
		return
	}
	this.Divide(l)
}

// 向量：单位化值 (0向量禁用)
func (this *Vec2[T]) Normalized() Vec2[T] {
	v := this.Clone()
	v.Normalize()
	return v
}

// 朝向不变拉长度
func (this *Vec2[T]) ScaleToLength(newLength T) {
	this.Normalize()
	this.Multiply(newLength)
}

// 复制朝向定长
func (this *Vec2[T]) ScaledToLength(newLength T) Vec2[T] {
	v := this.Clone()
	v.ScaleToLength(newLength)
	return v
}

// 返回：新向量
func NewVec2[T Numeric](x, y T) Vec2[T] {
	return Vec2[T]{X: x, Y: y}
}

// 返回：零向量(0,0,0)
func ZeroV2[T Numeric]() Vec2[T] {
	return Vec2[T]{X: 0, Y: 0}
}

// X 轴 单位向量
func XAxisV2[T Numeric]() Vec2[T] {
	return Vec2[T]{X: 1, Y: 0}
}

// Y 轴 单位向量
func YAxisV2[T Numeric]() Vec2[T] {
	return Vec2[T]{X: 0, Y: 1}
}

func XYAxisV2[T Numeric]() Vec2[T] {
	return Vec2[T]{X: 1, Y: 1}
}

func PositiveInfinityV2[T Numeric]() Vec2[T] {
	return Vec2[T]{X: T(math.Inf(1)), Y: T(math.Inf(1))}
}
func NegativeInfinityV2[T Numeric]() Vec2[T] {
	return Vec2[T]{X: T(math.Inf(-1)), Y: T(math.Inf(-1))}
}

// 返回：a + b 向量
func AddV2[T Numeric](a, b Vec2[T]) Vec2[T] {
	return Vec2[T]{X: a.X + b.X, Y: a.Y + b.Y}
}

// 返回：a - b 向量
func SubV2[T Numeric](a, b Vec2[T]) Vec2[T] {
	return Vec2[T]{X: a.X - b.X, Y: a.Y - b.Y}
}

func AddArrayV2[T Numeric](vs []Vec2[T], dv Vec2[T]) []Vec2[T] {
	for i, _ := range vs {
		vs[i].Add(dv)
	}
	return vs
}

func MultiplyV2[T Numeric](v Vec2[T], scalars []T) []Vec2[T] {
	var vs []Vec2[T]
	for _, value := range scalars {
		vector := v.Clone()
		vector.Multiply(value)
		vs = append(vs, vector)
	}
	return vs
}

// 求两点间距离
func DistanceV2[T Numeric](a Vec2[T], b Vec2[T]) T {
	return T(math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2)))
}

// 求两点间距离平方
func DistanceSqrV2[T Numeric](a Vec2[T], b Vec2[T]) T {
	return T(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}

// 线性插值
func LerpV2[T Numeric](a, b Vec2[T], t T) Vec2[T] {
	return Vec2[T]{X: a.X + (b.X-a.X)*t, Y: a.Y + (b.Y-a.Y)*t}
}

func LerpUnclampedV2[T Numeric](a, b Vec2[T], t T) Vec2[T] {
	return LerpV2(b, a, 1-t)
}

// 如果需要再加
//Max	返回由两个向量的最大分量组成的向量。
//Min	返回由两个向量的最小分量组成的向量。
//MoveTowards	将点 current 移向 /target/。
//Reflect	从法线定义的向量反射一个向量。
//Scale	将两个向量的分量相乘。
//SignedAngle	返回 from 与 to 之间的有符号角度（以度为单位）。
//SmoothDamp	随时间推移将一个向量逐渐改变为所需目标。
