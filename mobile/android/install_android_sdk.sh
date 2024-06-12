
# 定义 JDK 版本（8 或 11）
JDK_VERSION="11"

# 安装 OpenJDK
echo "Installing OpenJDK $JDK_VERSION..."
sudo apt update
sudo apt install -y openjdk-$JDK_VERSION-jdk-headless

# 设置 JAVA_HOME
export JAVA_HOME=/usr/lib/jvm/java-$JDK_VERSION-openjdk-amd64

# 下载 Android Command Line Tools
echo "Downloading Android Command Line Tools..."
cd ~
curl -L https://dl.google.com/android/repository/commandlinetools-linux-8512546_latest.zip -o /tmp/cmd-tools.zip
mkdir -p android/cmdline-tools
unzip -q -d android/cmdline-tools /tmp/cmd-tools.zip
mv android/cmdline-tools/cmdline-tools android/cmdline-tools/latest
rm /tmp/cmd-tools.zip

# 设置环境变量
echo "Setting up environment variables..."
echo "export ANDROID_HOME=\$HOME/android" >> ~/.bashrc
echo "export ANDROID_SDK_ROOT=\$ANDROID_HOME" >> ~/.bashrc
echo "export PATH=\$PATH:\$ANDROID_HOME/cmdline-tools/latest/bin:\$ANDROID_HOME/platform-tools:\$ANDROID_HOME/tools:\$ANDROID_HOME/tools/bin" >> ~/.bashrc
source ~/.bashrc

# 接受 SDK 许可
echo "Accepting SDK licenses..."
yes | ~/android/cmdline-tools/latest/bin/sdkmanager --licenses

# 更新并安装 SDK 组件
echo "Updating and installing SDK components..."
~/android/cmdline-tools/latest/bin/sdkmanager --update
~/android/cmdline-tools/latest/bin/sdkmanager "platforms;android-30" "build-tools;30.0.3"

echo "Android SDK installation completed."

echo "Android ndk install..."
sdkmanager "ndk-bundle"
echo "Android ndk installation completed."