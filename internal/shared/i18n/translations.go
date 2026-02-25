package i18n

import (
	"cms-api/internal/pkg/i18nutil"
)

var translations = map[Key]map[string]string{
	ErrNotFound: {
		i18nutil.LangEnglish: "The requested resource was not found",
		i18nutil.LangArabic:  "المورد المطلوب غير موجود",
	},
	ErrBadRequest: {
		i18nutil.LangEnglish: "Invalid request",
		i18nutil.LangArabic:  "طلب غير صالح",
	},
	ErrUnauthorized: {
		i18nutil.LangEnglish: "Authentication required",
		i18nutil.LangArabic:  "المصادقة مطلوبة",
	},
	ErrForbidden: {
		i18nutil.LangEnglish: "You don't have permission to perform this action",
		i18nutil.LangArabic:  "ليس لديك صلاحية للقيام بهذا الإجراء",
	},
	ErrConflict: {
		i18nutil.LangEnglish: "A conflict occurred with the current state",
		i18nutil.LangArabic:  "حدث تعارض مع الحالة الحالية",
	},
	ErrInternalServer: {
		i18nutil.LangEnglish: "An unexpected error occurred",
		i18nutil.LangArabic:  "حدث خطأ غير متوقع",
	},
	ErrServiceUnavailable: {
		i18nutil.LangEnglish: "Service is temporarily unavailable",
		i18nutil.LangArabic:  "الخدمة غير متاحة مؤقتاً",
	},
	ErrInvalidCredentials: {
		i18nutil.LangEnglish: "Invalid email or password",
		i18nutil.LangArabic:  "البريد الإلكتروني أو كلمة المرور غير صحيحة",
	},
	ErrInvalidToken: {
		i18nutil.LangEnglish: "Invalid or expired token",
		i18nutil.LangArabic:  "الرمز غير صالح أو منتهي الصلاحية",
	},
	ErrEmailAlreadyExists: {
		i18nutil.LangEnglish: "An account with this email already exists",
		i18nutil.LangArabic:  "يوجد حساب بهذا البريد الإلكتروني بالفعل",
	},
	ErrUserInactive: {
		i18nutil.LangEnglish: "Your account is inactive",
		i18nutil.LangArabic:  "حسابك غير نشط",
	},
	ErrTokenExpired: {
		i18nutil.LangEnglish: "Token has expired",
		i18nutil.LangArabic:  "انتهت صلاحية الرمز",
	},
	ErrTokenRevoked: {
		i18nutil.LangEnglish: "Token has been revoked",
		i18nutil.LangArabic:  "تم إلغاء الرمز",
	},
	ErrValidationFailed: {
		i18nutil.LangEnglish: "Validation failed",
		i18nutil.LangArabic:  "فشل التحقق من البيانات",
	},
}

func GetMessage(key Key, lang string) string {
	lang = i18nutil.NormalizeLanguage(lang)

	messages, exists := translations[key]
	if !exists {
		return string(key)
	}

	message, exists := messages[lang]
	if !exists {
		if englishMsg, ok := messages[i18nutil.LangEnglish]; ok {
			return englishMsg
		}
		return string(key)
	}

	return message
}
