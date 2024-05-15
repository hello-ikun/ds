class Node:
    def __init__(self, parent=None):
        self.parent = parent
        self.entries = []
        self.children = []

class Entry:
    def __init__(self, key, val):
        self.key = key
        self.val = val
        
class BTree:
    def __init__(self, order=3):
        if order < 3:
            order = 3
        self.root = None
        self.size = 0
        self.m = order

    def put(self, key, val):
        entry = Entry(key, val)
        # if root is None
        if self.root is None:
            self.root = Node()
            self.root.entries.append(entry)
            self.size += 1
            return
        # start to insert the data to tree
        if self._insert(self.root, entry):
            self.size += 1

    def _insert(self, node, entry):
        # if node is leaf
        # each insertion is on a leaf node
        if self._is_leaf(node):
            return self._insert_into_leaf(node, entry)
        # if node is not leaf to find leaf
        return self._insert_into_internal(node, entry)

    def _insert_into_leaf(self, node, entry):
        idx, ok = self._search(node, entry.key)
        if ok:
            node.entries[idx] = entry
            return False
        node.entries.insert(idx, entry)
        self._split(node)
        return True

    def _split(self, node):
        if not self._should_split(node):
            return
        print("Node:", node, "splitting")
        # root node split
        if node == self.root:
            self._split_root()
            return
        self._split_non_root(node)

    def _split_root(self):
        mid = self._middle()
        left = Node()
        left.entries = self.root.entries[:mid]
        right = Node()
        right.entries = self.root.entries[mid+1:]
        if not self._is_leaf(self.root):
            left.children = self.root.children[:mid+1]
            right.children = self.root.children[mid+1:]
            self._set_parent(left.children, left)
            self._set_parent(right.children, right)
        new_root = Node()
        new_root.entries = [self.root.entries[mid]]
        new_root.children = [left, right]
        left.parent = new_root
        right.parent = new_root
        self.root = new_root

    def _split_non_root(self, node):
        middle = self._middle()
        parent = node.parent
        left = Node(parent)
        left.entries = node.entries[:middle]
        right = Node(parent)
        right.entries = node.entries[middle+1:]
        if not self._is_leaf(node):
            left.children = node.children[:middle+1]
            right.children = node.children[middle+1:]
            self._set_parent(left.children, left)
            self._set_parent(right.children, right)
        insert_position, _ = self._search(parent, node.entries[middle].key)
        parent.entries.insert(insert_position, node.entries[middle])
        parent.children[insert_position] = left
        parent.children.insert(insert_position+1, right)
        # need to keep on passing up
        self._split(parent)

    def _set_parent(self, nodes, parent):
        for node in nodes:
            node.parent = parent

    def _should_split(self, node):
        return len(node.entries) > self._max_entries()

    def empty(self):
        return self.size == 0

    def _middle(self):
        return (self.m - 1) // 2

    def _max_children(self):
        return self.m

    def _min_children(self):
        return (self.m + 1) // 2

    def _max_entries(self):
        return self._max_children() - 1

    def _min_entries(self):
        return self._min_children() - 1

    def _is_leaf(self, node):
        return len(node.children) == 0
    
    # keep searching until find the final node.
    def _insert_into_internal(self, node, entry):
        idx, ok = self._search(node, entry.key)
        if ok:
            node.entries[idx] = entry
            return False
        return self._insert(node.children[idx], entry)
    # Implement a binary lookup and find where you need to insert it.
    def _search(self, node:Node, key):
        left, right = 0, len(node.entries) - 1
        ans = -1
        while left <= right:
            mid = (left + right) // 2
            if node.entries[mid].key == key:
                return mid, True
            elif node.entries[mid].key < key:
                ans = mid
                left = mid + 1
            else:
                right = mid - 1
        return ans + 1, False

        # use bisect_left to find key index
        # and notice data can sometimes not be found
        # the maximum length is returned ,but there is no problem
        # don't panic 
        # because the subtree is equal node +1
        # len(node.entries)+1==len(node.children)
        # idx=bisect.bisect_left(node.entries,key,key=lambda x:x.key)

        # return idx,idx<len(node.entries) and node.entries[idx].key==key
        

    def search_recur(self, start_node, key):
        # if node is empty return None
        if self.empty():
            return None, -1, False
        node = start_node
        while True:
            idx, ok = self._search(node, key)
            if ok:# if find the key in children return answer
                return node, idx, True
            # if not found and node is leaf return None
            if self._is_leaf(node):
                return None, -1, False
            # recur the data to get answer
            node = node.children[idx]

    def get(self, key):
        node, idx, ok = self.search_recur(self.root, key)
        if ok:
            return node.entries[idx].val, True
        return None, False
    # remove key from tree
    def remove(self, key):
        node, idx, ok = self.search_recur(self.root, key)
        if ok:
            self._delete(node, idx)
            self.size -= 1
    # delete the node base on index
    def _delete(self, node, idx):
        # if leaf node
        if self._is_leaf(node):
            del_key = node.entries[idx].key
            # del the leaf node and try to rebalance the tree
            del node.entries[idx]
            self._rebalance(node, del_key)
            if len(self.root.entries) == 0:
                self.root = None
            return
        # else not leaf node
        left_largest_node = self._right(node.children[idx])
        left_largest_entry_index = len(left_largest_node.entries) - 1
        node.entries[idx] = left_largest_node.entries[left_largest_entry_index]
        del_key = left_largest_node.entries[left_largest_entry_index].key
        del left_largest_node.entries[left_largest_entry_index]
        self._rebalance(left_largest_node, del_key)

    # delete data is node entries
    def _rebalance(self, node, delete_key):
        if node is None or len(node.entries) >= self._min_entries():
            return
        # find left sibling node
        left_sibling, left_sibling_index = self._left_sibling(node, delete_key)
        # if left sibling node fulfill a condition 
        # borrowing one from the left node is permissible
        if left_sibling is not None and len(left_sibling.entries) > self._min_entries():
            # if left sibling is exist and addition from the front position 
            node.entries.insert(0, node.parent.entries[left_sibling_index])
            node.parent.entries[left_sibling_index] = left_sibling.entries.pop()
            # deal with left sibling
            if not self._is_leaf(left_sibling):
                left_sibling_right_most_child = left_sibling.children.pop()
                left_sibling_right_most_child.parent = node
                node.children.insert(0, left_sibling_right_most_child)
            return
        # right sibling delete data and return index is next 
        right_sibling, right_sibling_index = self._right_sibling(node, delete_key)
        if right_sibling is not None and len(right_sibling.entries) > self._min_entries():
            # right sibling index -1 
            # because is parent entries
            node.entries.append(node.parent.entries[right_sibling_index - 1])
            node.parent.entries[right_sibling_index - 1] = right_sibling.entries.pop(0)
            # append right sibling left most child 
            if not self._is_leaf(right_sibling):
                right_sibling_left_most_child = right_sibling.children.pop(0)
                right_sibling_left_most_child.parent = node
                node.children.append(right_sibling_left_most_child)
            return
        
        # merge operation 
        if right_sibling is not None:
            # deal with node entries
            # node.entries.append(node.parent.entries[right_sibling_index - 1])
            # node.entries.extend(right_sibling.entries)
            
            # merge all the nodes.
            entries=node.entries+[node.parent.entries[right_sibling_index - 1]]+right_sibling.entries
            # delete node parent
            del_key = node.parent.entries[right_sibling_index - 1].key
            del node.parent.entries[right_sibling_index - 1]
            # set children parent
            self._append_children(node.parent.children[right_sibling_index], node)
            del node.parent.children[right_sibling_index]

            node.entries = entries

        elif left_sibling is not None:
            # all entries
            entries = left_sibling.entries + [node.parent.entries[left_sibling_index]] + node.entries
            # del the new key 
            del_key = node.parent.entries[left_sibling_index].key
            del node.parent.entries[left_sibling_index]

            self._prepend_children(node.parent.children[left_sibling_index], node)
            del node.parent.children[left_sibling_index]
            # set node entries
            node.entries = entries
        
        # if delete node parent is root
        if node.parent == self.root and len(self.root.entries) == 0:
            self.root = node
            node.parent = None
            return
        
        # rebalance parent node and update the delete key
        self._rebalance(node.parent, del_key)

    # seacch left node
    def _left(self, node):
        if self.empty():
            return None
        cur = node
        while True:
            if self._is_leaf(cur):
                return cur
            cur = cur.children[0]
    # seacch right node
    def _right(self, node):
        if self.empty():
            return None
        cur = node
        while True:
            if self._is_leaf(cur):
                return cur
            cur = cur.children[-1]
    # seacch left sibling node 
    # index -1
    def _left_sibling(self, node, key):
        if node.parent is not None:
            idx, _ = self._search(node.parent, key)
            idx -= 1
            if 0 <= idx < len(node.parent.children):
                return node.parent.children[idx], idx
        return None, -1
    # seacch right sibling node
    # index + 1
    def _right_sibling(self, node, key):
        if node.parent is not None:
            idx, _ = self._search(node.parent, key)
            idx += 1
            if idx < len(node.parent.children):
                return node.parent.children[idx], idx
        return None, -1
    # delete the child
    def _delete_child(self, node, idx):
        if idx >= len(node.children):
            return
        del node.children[idx]
    
    # append children
    # merge children node and set parent to_node
    def _append_children(self, from_node, to_node):
        to_node.children.extend(from_node.children)
        self._set_parent(from_node.children, to_node)
    # prepend children
    # merge all child nodes and set the father node
    def _prepend_children(self, from_node, to_node):
        to_node.children = from_node.children + to_node.children
        self._set_parent(from_node.children, to_node)
    def print(self):
        self._print_tree(self.root)

    def _print_tree(self, node, is_last=False, prefix=""):
        print(prefix + ("└── " if is_last else "├── ") + node.entries)

        child_count = len(node.children)
        for i, child in enumerate(node.children):
            is_last_child = i == child_count - 1
            self._print_tree(child, is_last=is_last_child, prefix=prefix + ("    " if is_last else "│   "))
    def print(self):
        self._print_tree(self.root)
    def _print_tree(self, node):
        print("\n******************************")

        def dsc(node):
            s = ""
            for keyword in node.entries:
                s += f"{keyword.key},"
            return s[:-1]

        logger = Logger()
        logger.tree(node, "child_nodes", dsc, 0)

        print("******************************")
class Logger:
    def tree(self, node, child_name, dsc, depth, prefix="    "):
        if depth == 0:
            print("bt--", dsc(node))
            depth += 1

        child_count = len(node.children)
        for idx, child in enumerate(node.children):
            is_last_child = idx == child_count - 1

            # if the current node has children, it's a branch node
            new_prefix = prefix + ("    " if is_last_child else "|   ")
            print(prefix + ("└-- " if is_last_child else "|-- ") + dsc(child))

            # recurse with the updated prefix
            self.tree(child, child_name, dsc, depth + 1, new_prefix)

if __name__ == "__main__":
    btree = BTree()
    for i in [20, 30, 50, 52, 60, 62]:
        btree.put(i,"val"+str(i))
    
    btree.print()
    btree.put(64,"val"+str(64))
    btree.print()
    btree.remove(52)
    btree.print()

    for i in [20, 30, 50, 52, 60, 62]:
        print(i,btree.get(i))

