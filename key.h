#pragma once
#include <string>
#include <typeinfo>
#include <iostream>
#include <variant>

using KeyType = std::variant<int, char, std::string>;

class Key {
	KeyType key;

	public:
		Key() = default;


		template <typename T>
		Key(T val) : key(val) {}
		//Key(T&& val) : key(std::forward<T>(val)) {}


		bool isInt() const;
		bool isChar() const;
		bool isString() const;
		const type_info& getType() const;
		std::string getTypeName() const;
		auto getValue() const -> decltype(key);

		friend std::ostream& operator<<(std::ostream& os, const Key& obj);
		friend bool operator<(const Key& lhs, const Key& rhs);
		friend bool operator>(const Key& lhs, const Key& rhs);
		friend bool operator==(const Key& lhs, const Key& rhs);
};

inline std::ostream& operator<<(std::ostream& os, const Key& obj) {
	if (obj.isInt()) {
		os << std::get<int>(obj.key);
	} 
	else if (obj.isChar()) {
		os << std::get<char>(obj.key);
	}
	else if (obj.isString()) {
		os << std::get<std::string>(obj.key);
	}
	return os;
}

inline bool operator<(const Key& lhs, const Key& rhs) {
	if (lhs.isInt() && rhs.isInt()) {
		return std::get<int>(lhs.key) < std::get<int>(rhs.key);
	}
	else if (lhs.isChar() && rhs.isChar()) {
		return std::get<char>(lhs.key) < std::get<char>(rhs.key);
	}
	else if (lhs.isString() && rhs.isString()) {
		return std::get<std::string>(lhs.key) < std::get<std::string>(rhs.key);
	}
	else {
		throw std::runtime_error("Operator \"<\" Error: Invalid comparison");
	}
}

inline bool operator>(const Key& lhs, const Key& rhs) {
	if (lhs.isInt() && rhs.isInt()) {
		return std::get<int>(lhs.key) > std::get<int>(rhs.key);
	}
	else if (lhs.isChar() && rhs.isChar()) {
		return std::get<char>(lhs.key) > std::get<char>(rhs.key);
	}
	else if (lhs.isString() && rhs.isString()) {
		return std::get<std::string>(lhs.key) > std::get<std::string>(rhs.key);
	}
	else {
		throw std::runtime_error("Operator \">\" Error: Invalid comparison");
	}
}

inline bool operator==(const Key& lhs, const Key& rhs) {
	if (lhs.isInt() && rhs.isInt()) {
		return std::get<int>(lhs.key) == std::get<int>(rhs.key);
	}
	else if (lhs.isChar() && rhs.isChar()) {
		return std::get<char>(lhs.key) == std::get<char>(rhs.key);
	}
	else if (lhs.isString() && rhs.isString()) {
		return std::get<std::string>(lhs.key) == std::get<std::string>(rhs.key);
	}
	else {
		throw std::runtime_error("Operator \"==\" Error: Invalid comparison");
	}
}