import sys

def main():
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    # 尝试不同编码读取
    for encoding in ['utf-16', 'utf-8', 'utf-8-sig', 'gbk']:
        try:
            with open(input_file, 'r', encoding=encoding) as f:
                content = f.read()
            if 'COPY' in content:
                print(f"Read with encoding: {encoding}")
                break
        except:
            continue
    
    # 找到COPY数据块
    copy_start = content.find('COPY public.companies')
    if copy_start == -1:
        print("ERROR: Cannot find COPY statement")
        return
    
    # 找到数据开始位置
    data_start = content.find('\n', copy_start) + 1
    # 找到数据结束位置
    data_end = content.find('\n\\.', data_start)
    if data_end == -1:
        data_end = len(content)
    
    data_block = content[data_start:data_end].strip()
    lines = data_block.split('\n')
    
    inserts = []
    for line in lines:
        if not line.strip():
            continue
        
        # 用 tab 分割
        values = line.split('\t')
        if len(values) < 10:
            continue
        
        # 提取字段并处理特殊字符
        def clean(val):
            val = val.strip()
            if val == '\\N' or not val:
                return ''
            return val.replace("'", "''")
        
        id_val = values[0].strip()
        name = clean(values[1])
        name_en = clean(values[2])
        logo_url = clean(values[3])
        website = clean(values[4])
        industry = clean(values[5])
        group = clean(values[6]) or 'bigtech'
        description = clean(values[7])
        is_preset = values[8].strip()
        health_status = clean(values[9])[:10]  # 限制10个字符
        
        # 处理布尔值
        is_preset = 'true' if is_preset in ('t', 'true', '1') else 'false'
        
        # 处理ID
        if id_val == '\\N' or not id_val:
            id_val = 'gen_random_uuid()'
        else:
            id_val = f"'{id_val}'"
        
        # 生成INSERT语句
        insert = f"""INSERT INTO companies (id, name, name_en, logo_url, website, industry, "group", description, is_preset, health_status) VALUES ({id_val}, '{name}', '{name_en}', '{logo_url}', '{website}', '{industry}', '{group}', '{description}', {is_preset}, '{health_status}');"""
        inserts.append(insert)
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write("BEGIN;\n\n")
        for insert in inserts:
            f.write(insert + "\n")
        f.write("\nCOMMIT;\n")
    
    print(f"Generated {len(inserts)} INSERT statements")

if __name__ == '__main__':
    main()