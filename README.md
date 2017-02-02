media
=====
该组件用于分析多媒体文件相关属性，比如图片的类型/大小/尺寸，视频的类型/分辨率/时长等

## video
借用ffmpeg工具集来获取视频文件相关属性，等同于ffmpeg -i video输出

依赖ffmpeg工具集，并将其bin路径添加到PATH环境中, 保证命令ffmpeg能直接调用

## image
能够识别gif jpeg png bmp riff tiff webp等格式图片

依赖image和golang.org/x/image
