# coding: utf-8
from sqlite3 import connect


def create_db(con: connect) -> None:
    cur = con.cursor()
    cur.execute("""
                create table if not exists author(id    integer primary key autoincrement,
                                                  name  text,
                                                  descr text);
                """)
    cur.execute("""
                create table if not exists item(id          integer primary key autoincrement,
                                                id_author   integer references author(id) on delete cascade,
                                                name        text,
                                                next        integer);
                """)
    cur.execute("""
                create table if not exists line(id      integer primary key autoincrement,
                                                id_item integer references item(id) on delete cascade,
                                                num     integer
                                                ttext   text,
                                                stop    integer);
                """)

    cur.execute("""
                create table if not exists trans(id     integer primary key autoincrement,
                                                id_line integer references line(id) on delete cascade,
                                                ttext   text,
                                                lang    text);
                """)