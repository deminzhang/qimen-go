# 奇门遁甲 转盘/飞盘/鸣法排盘

# 基础奇门起局
1. 鸣法以满盘转时干为暗干
2. 鸣法九星只顺不逆
3. 鸣法锁定拆补法
4. 时家带大六壬的天地盘
5. 置闰法按:大雪芒种,年十位数奇数,个位置0, 3, 5, 8;十位数偶数,个位置1, 4, 7置闰,
 符头跨小雪/小满接气,上元跨二至超神
6. NASA数据存缓存到本地sqlite减少重复请求
7. 附星盘星图
8. 附八字排盘


# TODO
1. 查第三方节气精确时间跟天文授时误差问题
2. 日家奇门2:太乙九星派+黄黑道
3. 美化字色
4. UI坐标承父坐标
5. 八字的流通数据模拟
6. 星盘校正七政四余星宿角度
7. 星盘计算合拱冲刑

# 总结
1. 建星的位置就是太阳的位置,月将初入点
2. 木星,岁星就是年柱的力量


# 引用
	github.com/6tail/lunar-go //年历
	github.com/hajimehoshi/ebiten/v2 //2D游戏引擎

# build apk
```shell
ebitenmobile bind -target android -javapkg com.deminzhang.qimen -o ./mobile/android/qimen/qimen.aar ./mobile
cd ./mobile/android
./gradlew assembleRelease