# coding: utf-8

import os
import db


def ins():
    # get author from user
    author = db.get_author(author_name=input("Author: "))

    # get item from user
    item = db.get_item(author=author, item_name=input("Item: "))
    num = 0

    f_name = item.name + '.poem'

    with open(f_name, 'r') as f:
        for line in f.readlines():
            num += 1
            item.l_lines.append(db.Line(tuple([0, item.id, num, line.strip()])))

    if not int(author.exists): author.ins()
    if not int(item.exists): 
        item.ins()
        for line in item.l_lines: 
            line.ins()
    
    os.remove(f_name)

if __name__ == "__main__":
    ins()