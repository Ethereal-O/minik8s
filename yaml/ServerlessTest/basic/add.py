def add(d):
    x = d.get('x', 0)
    y = d.get('y', 0)
    result = int(x) + int(y)
    d['result'] = str(result)
    return d