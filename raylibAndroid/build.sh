
# export ANDROID_NDK_HOME=/opt/android-ndk
# export ANDROID_HOME=/opt/android-sdk
go clean
rm android/libs/arm64-v8a/libexample.so
export PATH=${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/bin:${PATH}
export ANDROID_SYSROOT=${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/sysroot
export ANDROID_API=26
export ANDROID_TOOLCHAIN=${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64
CC="aarch64-linux-android${ANDROID_API}-clang" CGO_CFLAGS="-I${ANDROID_SYSROOT}/usr/include -I${ANDROID_SYSROOT}/usr/include/aarch64-linux-android --sysroot=${ANDROID_SYSROOT}" CGO_LDFLAGS="-L${ANDROID_SYSROOT}/usr/lib/aarch64-linux-android/${ANDROID_API} \
-L${ANDROID_TOOLCHAIN}/aarch64-linux-android/lib --sysroot=${ANDROID_SYSROOT}" \
CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
go build -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-O5,-soname,libexample.so" -o=android/libs/arm64-v8a/libexample.so
${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt/linux-x86_64/bin/llvm-strip android/libs/arm64-v8a/libexample.so
chmod +x android/libs/arm64-v8a/libexample.so
./upx --ultra-brute --android-shlib android/libs/arm64-v8a/libexample.so
./gradlew assembleDebug