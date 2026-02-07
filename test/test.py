# Python 测试文件

# 全局变量
global_config = "config"
DEBUG_MODE = True
cache_data = {}

def large_function():
    """这是一个超过200行的大函数"""
    print("Line 1")
    print("Line 2")
    print("Line 3")
    print("Line 4")
    print("Line 5")
    # ... 省略中间部分，实际应该有200+行
    for i in range(100):
        print(f"Processing {i}")
        if i % 2 == 0:
            print("Even")
        else:
            print("Odd")

    for j in range(100):
        print(f"Another loop {j}")
        if j % 3 == 0:
            print("Divisible by 3")

    return "Done"

def small_function():
    """这是一个小函数"""
    return "Small"
