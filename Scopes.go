package VkApi

//	Пользователь разрешил отправлять ему уведомления (для flash/iframe-приложений).
const U_NOTIFY = 1

//	Доступ к друзьям.
const U_FRIENDS = 2

//	Доступ к фотографиям.
const U_PHOTOS = 4

//	Доступ к аудиозаписям.
const U_AUDIO = 8

//	Доступ к видеозаписям.
const U_VIDEO = 16

//	Доступ к историям.
const U_STORIES = 64

//	Доступ к wiki-страницам
const U_PAGES = 128

//	Добавление ссылки на приложение в меню слева.
const U_LINK = 256

//	Доступ к статусу пользователя.
const U_STATUS = 1024

//	Доступ к заметкам пользователя.
const U_NOTES = 2048

//	Доступ к расширенным методам работы с сообщениями (только для Standalone-приложений).
const U_MESSAGES = 4096

//	Доступ к обычным и расширенным методам работы со стеной.
// Данное право доступа по умолчанию недоступно для сайтов (игнорируется при попытке авторизации для приложений с типом «Веб-сайт» или по схеме Authorization Code Flow (https://vk.com/dev/authcode_flow_user) ).
const U_WALL = 8192

//	Доступ к расширенным методам работы с рекламным API (https://vk.com/dev/ads). Доступно для авторизации по схеме Implicit Flow (https://vk.com/dev/implicit_flow_user) или Authorization Code Flow (https://vk.com/dev/authcode_flow_user).
const U_ADS = 32768

//	Доступ к API (https://vk.com/dev/apiusage) в любое время (при использовании этой опции параметр expires_in, возвращаемый вместе с access_token, содержит 0 — токен бессрочный). Не применяется в Open API.
const U_OFFLINE = 65536

//	Доступ к документам.
const U_DOCS = 131072

//	Доступ к группам пользователя.
const U_GROUPS = 262144

//	Доступ к оповещениям об ответах пользователю.
const U_NOTIFICATIONS = 524288

//	Доступ к статистике групп и приложений пользователя, администратором которых он является.
const U_STATS = 1048576

//	Доступ к email пользователя.
const U_EMAIL = 4194304

//	Доступ к товарам.
const U_MARKET = 134217728

//	Доступ к историям.
const G_STORIES = 1

//	Доступ к фотографиям.
const G_PHOTOS = 4

//	Доступ к виджетам приложений сообществ. Это право можно запросить только с помощью метода Client API showGroupSettingsBox.
const G_APP_WIDGET = 64

//	Доступ к сообщениям сообщества.
const G_MESSAGES = 4096

//	Доступ к документам.
const G_DOCS = 131072

//	Доступ к администрированию сообщества.
const G_MANAGE = 262144
