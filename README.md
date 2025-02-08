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
9. 附梅花易数时盘,手动输上下变数可起


# TODO
1. 查第三方节气精确时间跟天文授时误差问题
2. 日家奇门2:太乙九星派+黄黑道
3. 八字的流通数据模拟
4. 星盘校正七政四余星宿角度
5. 星盘计算合拱冲刑
6. 大六壬排四课三传
7. NASA数据缓存优化为只记整点,每分数据用插值近似计算
8. 阴盘奇门排盘
9. 常用问事简解,记录对比正确率

# 总结心得
1. 建星的位置就是太阳的位置,月将为初入点,月建为中点
2. 木星,岁星就是年柱的力量
3. "见实相，诸法空，刹那顿悟万法同"
4. 陶氏,鸣法奇门的格局同大六壬,可能是以大六壬思路推导的
5. 考阴盘奇门应为洪范奇门传入韩国后以梅花易数起局方式改编
6. 大六壬的的起三传总结为寻找局中唯一,最初,最重的矛盾源头点定初传,定局中主要力量源


# 第三方引用
	github.com/6tail/lunar-go //万年历 （MIT License）
	github.com/hajimehoshi/ebiten/v2 //2D游戏引擎 （Apache License 2.0）

# build apk
```shell
ebitenmobile bind -target android -javapkg com.deminzhang.qimen -o ./mobile/android/qimen/qimen.aar ./mobile
cd ./mobile/android
./gradlew assembleRelease