#pragma once
#include <string>
#include <typeinfo>
#include <variant>
#include <sstream>

class Key {
	std::variant<int, char, std::string> value;

	public:
		Key() {}

		template <typename T>
		Key(const T& val) : value(val) {}

		const type_info& getType();
		std::string getTypeName();
		auto getValue() -> decltype(value);

		friend std::ostream& operator<<(std::ostream& os, const Key& obj);
		friend bool operator<(const Key& lhs, const Key& rhs);
		friend bool operator>(const Key& lhs, const Key& rhs);
		friend bool operator==(const Key& lhs, const Key& rhs);
	private:
		bool isInt() const;
		bool isChar() const;
		bool isString() const;
};

inline std::ostream& operator<<(std::ostream& os, const Key& obj) {
	std::visit([&os](const auto& val) { os << val; }, obj.value);
	return os;
}

inline bool operator<(const Key& lhs, const Key& rhs) {
	return std::visit([](const auto& lv, const auto& rv) {
		if constexpr (std::is_same_v<decltype(lv), decltype(rv)>) {
			return lv < rv;
		}
		else {
			printf("Operator \"<\" Warning: Types are different\n");
			return false; 
		}
		}, lhs.value, rhs.value);
}

inline bool operator>(const Key& lhs, const Key& rhs) {
	return std::visit([](const auto& lv, const auto& rv) {
		if constexpr (std::is_same_v<decltype(lv), decltype(rv)>) {
			return lv > rv;
		}
		else {
			printf("Operator \">\" Warning: Types are different\n");
			return false; 
		}
		}, lhs.value, rhs.value);
}

inline bool operator==(const Key& lhs, const Key& rhs) {
	return std::visit([](const auto& lv, const auto& rv) {
		if constexpr (std::is_same_v<decltype(lv), decltype(rv)>) {
			return lv == rv;
		}
		else {
			printf("Operator \"==\" Warning: Types are different\n");
			return false; 
		}
		}, lhs.value, rhs.value);
}

struct KeyComparator {
	bool operator()(const Key& lhs, const Key& rhs) const {
		return lhs < rhs; 
	}
};