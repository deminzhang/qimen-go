# 奇门遁甲 转盘/飞盘/鸣法排盘

# 基础奇门起局

1. 鸣法以满盘转时干为暗干
2. 鸣法九星只顺不逆
3. 鸣法锁定拆补法
4. 时家带大六壬的月将支盘以查天门地户


# TODO
1. 查第三方节气精确时间跟天文授时误差问题
2. 起局置闰法
3. 日家奇门1:天蓬九星派
4. 日家奇门2:太乙九星派+黄黑道
5. 美化字色
6. UI坐标承父坐标
7. 阴盘奇门 年/时/刻


# build apk
```shell
ebitenmobile bind -target android -javapkg com.deminzhang.qimen -o ./mobile/android/qimen/qimen.aar ./mobile
cd ./mobile/android
./gradlew assembleRelease