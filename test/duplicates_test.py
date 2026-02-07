# Python 重复代码测试文件

def process_user(name, age):
    """处理数据"""
    print(f"Processing: {name}, age: {age}")
    if age < 18:
        print("Minor")
    else:
        print("Adult")

def process_customer(name, age):
    """处理数据"""
    print(f"Processing: {name}, age: {age}")
    if age < 18:
        print("Minor")
    else:
        print("Adult")

def validate_email(email):
    """检查邮箱"""
    if len(email) == 0:
        return False
    return True

def check_email(email):
    """检查邮箱"""
    if len(email) == 0:
        return False
    return True

def unique_function():
    """这是一个独特的函数"""
    print("This is unique")
    print("Very unique")
