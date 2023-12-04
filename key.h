#pragma once
#include <string>
#include <typeinfo>
#include <variant>
#include <sstream>

using Value = std::variant<int, double, char, std::string>;

class Object {
	Value value;

	public:
		template <typename T>
		Object(T val) : value(val) {}

		const type_info& getType();
		std::string getTypeName();
		auto getValue() -> decltype(value);

		friend std::ostream& operator<<(std::ostream& os, const Object& obj);

	private:
		bool isInt() const;
		bool isDouble() const;
		bool isChar() const;
		bool isString() const;
};

inline
std::ostream& operator<<(std::ostream& os, const Object& obj) {
	std::visit([&os](const auto& val) { os << val; }, obj.value);
	return os;
}

