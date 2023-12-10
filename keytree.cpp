/*
KeyNode* KeyTree::insert(Key *key) {
    if (this->root == nullptr) {
		this->root = new KeyNode(key);
		return this->root;
	}
	else {
		KeyNode* current = this->root;
		while (true) {
			if (key < current->key) {
				if (current->children.size() == 0) {
					KeyNode* newNode = new KeyNode(key);
					newNode->parent = current;
					current->children.push_back(newNode);
					return newNode;
				}
				else {
					current = current->children[0];
				}
			}
			else if (current->key < key) {
				if (current->children.size() == 0) {
					KeyNode* newNode = new KeyNode(key);
					newNode->parent = current;
					current->children.push_back(newNode);
					return newNode;
				}
				else {
					current = current->children[current->children.size() - 1];
				}
			}
			else {
				return current;
			}
		}
	}
}

KeyNode* KeyTree::search(Key *key) {
	if (this->root == nullptr) {
		return nullptr;
	}

	KeyNode* current = this->root;
	
	while (true) {
		if (key < current->key) {
			if (current->children.size() == 0) {
				return nullptr;
			}
			else {
				current = current->children[0];
			}
		} else
		if (current->key < key) {
			if (current->children.size() == 0) {
				return nullptr;
			}
			else {
				current = current->children[current->children.size() - 1];
			}
		}
		else {
			return current;
		}
	}
	
}

void KeyTree::remove(Key *key) {
	if (this->root == nullptr) {
		return;
	}

	KeyNode* current = this->search(key);

	if (current == nullptr) {
		return;
	}

	if (current->children.size() == 0) {
		if (current->parent == nullptr) {
			this->root = nullptr;
		}
		else {
			for (int i = 0; i < current->parent->children.size(); i++) {
				if (current->parent->children[i] == current) {
					current->parent->children.erase(current->parent->children.begin() + i);
					break;
				}
			}
		}
	}
	else {
		KeyNode* replacement = current->children[0];
		while (replacement->children.size() > 0) {
			replacement = replacement->children[replacement->children.size() - 1];
		}

		current->key = replacement->key;

		if (replacement->parent == nullptr) {
			current->children.erase(current->children.begin());
		}
		else {
			for (int i = 0; i < replacement->parent->children.size(); i++) {
				if (replacement->parent->children[i] == replacement) {
					replacement->parent->children.erase(replacement->parent->children.begin() + i);
					break;
				}
			}
		}
	}	
}*/