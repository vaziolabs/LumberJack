#include "keytree.h"

KeyNode* KeyTree::insert(KeyNode* keynode) {
	if (this->root == nullptr) {
		this->root = keynode;
		return keynode;
	}

	KeyNode* current = this->root;

	while (current != nullptr) {
		if (keynode->key < current->key) {
			if (current->children.size() == 0) {
				current->children.push_back(keynode);
				keynode->parent = current;
				return keynode;
			}
			else {
				current = current->children[0];
			}
		}
		else if (keynode->key > current->key) {
			if (current->children.size() == 0) {
				current->children.push_back(keynode);
				keynode->parent = current;
				return keynode;
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

KeyNode* KeyTree::search(KeyType key) const {
	if (this->root == nullptr) {
		return nullptr;
	}

	KeyNode* current = this->root;

	while (current != nullptr) {
		if (key < current->key) {
			if (current->children.size() == 0) {
				return nullptr;
			}
			else {
				current = current->children[0];
			}
		}
		else if (key > current->key) {
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

void KeyTree::remove(int id) const {
	if (this->root == nullptr) {
		return;
	}

	KeyNode* current = this->root;

	if (!current->hasDescendant(id)) { return; }

	while (current != nullptr) {
		if (id < current->key) {
			if (current->children.size() == 0) { return; } 
			else { current = current->children[0]; }
		} else if (id > current->key) {
			if (current->children.size() == 0) { return; }
			else { current = current->children[current->children.size() - 1]; }
		} else {
			if (current->children.size() == 0) {
				if (current->parent->children.size() == 1) {
					current->parent->children.clear();
				} else {
					for (int i = 0; i < current->parent->children.size(); i++) {
						if (current->parent->children[i] == current) {
							current->parent->children.erase(current->parent->children.begin() + i);
							break;
						}
					}
				}
			} else if (current->children.size() == 1) {
				if (current->parent->children.size() == 1) {
					current->parent->children.clear();
				} else {
					for (int i = 0; i < current->parent->children.size(); i++) {
						if (current->parent->children[i] == current) {
							current->parent->children.erase(current->parent->children.begin() + i);
							break;
						}
					}
				}
			} else {
				KeyNode* replacement = current->children[0];

				while (replacement->children.size() > 0) {
					replacement = replacement->children[replacement->children.size() - 1];
				}

				current->key = replacement->key;

				if (replacement->parent->children.size() == 1) {
					replacement->parent->children.clear();
				} else {
					for (int i = 0; i < replacement->parent->children.size(); i++) {
						if (replacement->parent->children[i] == replacement) {
							replacement->parent->children.erase(replacement->parent->children.begin() + i);
							break;
						}
					}
				}
			}
		}
	}
}

void KeyTree::remove(char id) const { 
	if (this->root == nullptr) {
		return;
	}

	KeyNode* current = this->root;

	if (!current->hasDescendant(id)) { return; }

	while (current != nullptr) {
		if (id < current->key) {
			if (current->children.size() == 0) { return; } 
			else { current = current->children[0]; }
		} else if (id > current->key) {
			if (current->children.size() == 0) { return; }
			else { current = current->children[current->children.size() - 1]; }
		} else {
			if (current->children.size() == 0) {
				if (current->parent->children.size() == 1) {
					current->parent->children.clear();
				}
				else {
					for (int i = 0; i < current->parent->children.size(); i++) {
						if (current->parent->children[i] == current) {
							current->parent->children.erase(current->parent->children.begin() + i);
							break;
						}
					}
				}
			} else if (current->children.size() == 1) {
				if (current->parent->children.size() == 1) {
					current->parent->children.clear();
				} else {
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

				if (replacement->parent->children.size() == 1) {
					replacement->parent->children.clear();
				} else {
					for (int i = 0; i < replacement->parent->children.size(); i++) {
						if (replacement->parent->children[i] == replacement) {
							replacement->children.erase(replacement->parent->children.begin() + i);
						}
					}
				}
			}
		}
	}
}

void KeyTree::remove(std::string id) const {
	if (this->root == nullptr) {
		return;
	}

	KeyNode* current = this->root;

	if (!current->hasDescendant(id)) { return; }
	while (current != nullptr) {
		if (id < current->key) {
			if (current->children.size() == 0) { return; } 
			else { current = current->children[0]; }
		} else if (id > current->key) {
			if (current->children.size() == 0) { return; }
			else { current = current->children[current->children.size() - 1]; }
		} else {
			if (current->children.size() == 0) {
				if (current->parent->children.size() == 1) {
					current->parent->children.clear();
				} else {
					for (int i = 0; i < current->parent->children.size(); i++) {
						if (current->parent->children[i] == current) {
							current->parent->children.erase(current->parent->children.begin() + i);
							break;
						}
					}
				}
			} else if (current->children.size() == 1) {
				if (current->parent->children.size() == 1) {
					current->parent->children.clear();
				} else {
					for (int i = 0; i < current->parent->children.size(); i++) {
						if (current->parent->children[i] == current) {
							current->parent->children.erase(current->parent->children.begin() + i);
							break;
						}
					}
				}
			} else {
				KeyNode* replacement = current->children[0];

				while (replacement->children.size() > 0) {
					replacement = replacement->children[replacement->children.size() - 1];
				}

				current->key = replacement->key;

				if (replacement->parent->children.size() == 1) {
					replacement->parent->children.clear();
				} else {
					for (int i = 0; i < replacement->parent->children.size(); i++) {
						if (replacement->parent->children[i] == replacement) {
							replacement->parent->children.erase(replacement->parent->children.begin() + i);
						}
					}
				}
			}
		}
	}
}

void KeyTree::print() const {
	if (this->root == nullptr) {
		return;
	}

	std::vector<KeyNode*> currentLevel;
	std::vector<KeyNode*> nextLevel;

	currentLevel.push_back(this->root);

	while (currentLevel.size() > 0) {
		for (int i = 0; i < currentLevel.size(); i++) {
			std::cout << currentLevel[i]->key << " ";
			
			for (int j = 0; j < currentLevel[i]->children.size(); j++) {
				nextLevel.push_back(currentLevel[i]->children[j]);
			}
		}

		std::cout << std::endl;

		currentLevel = nextLevel;
		nextLevel.clear();
	}
}	