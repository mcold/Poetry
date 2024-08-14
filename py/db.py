from sqlite3 import connect

from pathlib import Path
from os import sep

db = 'C:\\Users\\mholo\\WD\\GO\\pro\\Poetry\\DB.db'
# db = str(Path.home()) + sep + 'Zotero\\zotero.sqlite'
# TODO: get path from env PATH

class Line:

    def __init__(self, t: tuple):
        if len(t) > 0:
            self.id = t[0]
            self.id_item = t[1]
            self.num = t[2]
            self.ttext = t[3]
            if len(t) > 4:
                self.stop = t[4]
            else:
                self.stop = 0
    
    def ins(self):
        with connect(db) as conn:
            cur = conn.cursor()
            cur.execute("""insert into line(id_item, num, ttext, stop)
                                   values({id_item}, '{num}', '{ttext}', {stop})
                              """.format(id_item=self.id_item, 
                                         num = self.num,
                                         ttext=self.ttext, 
                                         stop=self.stop))
            conn.commit()

    def __str__(self) -> str:
        return """\nLine\n{sep}\nID item: {id_item}\nNum: {num}\nText: {text}\nStop: {stop}""".format(sep='-'*40,
                             id_item=self.id_item, 
                             num=self.num, 
                             text=self.ttext, 
                             stop=self.stop)

class Item:

    def __init__(self, t: tuple):
        if len(t) > 0:
            self.id = t[0]
            self.id_author = t[1]
            self.name = t[2]
            self.next = t[3] if t[3] else 0
            self.exists = t[4]
            self.l_lines = list()

    def ins(self):
        with connect(db) as conn:
            cur = conn.cursor()
            cur.execute("""insert into item(id, id_author, name, next)
                                   values({id}, {id_author}, '{name}', {next})
                              """.format(id=self.id,
                                         id_author=self.id_author,
                                         name=self.name, 
                                         next=self.next))

            conn.commit()

    def __str__(self) -> str:
        return """\nItem\n{sep}\nID item: {id_item}\nID author: {id_author}\nName author: {name_author}\nNext: {next}""".format(sep='-'*40, 
                             id_item=self.id, 
                             id_author=self.id_author, 
                             name_author=self.name, 
                             next=self.next)

class Author:

    l_items = list()

    def __init__(self, t: tuple):
        if len(t) > 0:
            self.id = t[0]
            self.name = t[1]
            self.descr = t[2] if t[2] else ''
            self.exists = t[3]
            self.l_items = list()
        else:
            self.id = None
            self.name = None
            self.descr = None
            self.exists = 0
            self.l_items = list()
    
    def add_item(self, item: Item):
        self.l_items.append(item)

    def ins(self):
        with connect(db) as conn:
            cur = conn.cursor()
            cur.execute("""insert into author(id, name, descr)
                                   values({id}, '{name}', '{descr}')
                              """.format(id=self.id, 
                                         name=self.name, 
                                         descr=self.descr))
            conn.commit()

    
    def __str__(self) -> str:
        return """\nAuthor\n{sep}\nID author: {id_author}\nName author: {name_author}\n""".format(sep='-'*40, 
                            id_author=self.id,
                            name_author=self.name)


def get_author(author_name: str):
    with connect(db) as conn:
        cur = conn.cursor()
        cur.execute("""select id,
                              name,
                              descr,
                              1
                         from author
                        where lower(name) = lower('{author_name}');
            """.format(author_name=author_name))
        try:
            return Author(cur.fetchone())
        except:
            cur.execute("""select max(id) + 1
                             from author""")
            return Author(tuple([cur.fetchone()[0], author_name, '', 0]))

def get_item(author: Author, item_name: str):
    with connect(db) as conn:
        cur = conn.cursor()
        cur.execute("""select id,
                              id_author,
                              name,
                              next,
                              1
                         from item
                        where id_author = {id_author}
                          and lower(name) = lower('{item_name}');
            """.format(id_author=author.id, item_name=item_name))
        try:
            return Item(cur.fetchone())
        except:
            cur.execute("""select max(id) + 1
                             from item""")
            return Item(tuple([cur.fetchone()[0], author.id, item_name, 0, 0]))


