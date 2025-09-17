package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Order struct { // вот наша новая структура
	ID          int    // поля структур, которые передаются в шаблон
	Title       string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	Price       string
	Value       string
	Img         string
	Description string
	Source      string
}

func (r *Repository) GetOrders() ([]Order, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	orders := []Order{ // массив элементов из наших структур
		{
			ID:          1,
			Title:       "Кирпич строительный (Российская Империя)",
			Price:       "2 руб. 50 коп.",
			Value:       "1000 шт.",
			Img:         "http://localhost:9000/iu5-web/img/object_1.jpg",
			Description: "Стандартный полнотелый кирпич ручной формовки, производившийся на московских заводах в конце XIX века. Основной материал для капитального строительства доходных домов, промышленных зданий и общественных учреждений. Цена указана с доставкой на стройплощадку в пределах города.",
			Source:      "«Справочник московского архитектора и инженера», 1895 г.",
		},
		{
			ID:          2,
			Title:       "Труд каменщика (Викторианская Англия)",
			Price:       "2 шиллингов",
			Value:       "10 часов",
			Img:         "http://localhost:9000/iu5-web/img/object_2.jpg",
			Description: "Работа высококвалифицированного каменщика, занятого на возведении фасадов и несущих конструкций общественных зданий. В оплату включено пользование собственным набором инструментов. Рабочий день длился 10 часов с одним перерывом.",
			Source:      "Отчет Британского Министерства труда по заработной плате в строительстве, 1850 г.",
		},
		{
			ID:          3,
			Title:       "Мрамор пентелийский (Древняя Греция)",
			Price:       "180 драхм",
			Value:       "1 талант веса",
			Img:         "http://localhost:9000/iu5-web/img/object_3.jpg",
			Description: "Высококачественный белый мрамор с легким золотистым оттенком, добываемый в каменоломнях горы Пенделикон к северо-востоку от Афин. Был основным материалом для строительства Афинского акрополя. Цена включает добычу и черновую обработку блоков на месте.",
			Source:      "Строительные надписи (Афинский акрополь, учетные записи сметы Парфенона), ок. 434-432 гг. до н.э.",
		},
		{
			ID:          4,
			Title:       "Свинец кровельный (Средневековая Франция)",
			Price:       "10 денье",
			Value:       "лист",
			Img:         "http://localhost:9000/iu5-web/img/object_4.jpg",
			Description: "Листовой свинец, отлитый в стандартные пластины размером приблизительно 2 x 1 фут. Использовался для покрытия крыш и желобов готических соборов, обеспечивая водонепроницаемость. Цена указана за лист без учета сложных работ по монтажу.",
			Source:      "Счетные книги Собора Нотр-Дам де Пари (фрагменты), XIII век.",
		},
		{
			ID:          5,
			Title:       "Масло льняное для живописи (Нидерланды, Золотой век)",
			Price:       "5 стюверов",
			Value:       "унцию",
			Img:         "http://localhost:9000/iu5-web/img/object_5.png",
			Description: "Очищенное льняное масло высшего качества, используемое в качестве связующего вещества для изготовления масляных красок. Продавалось в небольших пузырьках в специализированных лавках для художников и аптекарей.",
			Source:      "Опись имущества и счетные книги мастерской Рембрандта ван Рейна, Амстердам, 1656 г.",
		},
		{
			ID:          6,
			Title:       "Труд разнорабочего (США, Великая Депрессия)",
			Price:       "0.50 $",
			Value:       "час",
			Img:         "http://localhost:9000/iu5-web/img/object_6.jpg",
			Description: "Неквалифицированный физический труд на строительных работах: земляные работы, перенос материалов, подсобные операции. Работа часто была временной и низкооплачиваемой из-за огромного предложения рабочей силы в период экономического кризиса.",
			Source:      "Статистика заработной платы Бюро трудовой статистики США (U.S. BLS), 1932 г.",
		},
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("Array is empty")
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}

	return result, nil
}

type Estimate struct {
	ID       int
	OrderIDs []int
}

var estimates = []Estimate{
	{
		ID:       1,
		OrderIDs: []int{1, 2},
	},
}

func (r *Repository) GetEstimateData(id int) (Estimate, error) {
	for _, estimate := range estimates {
		if estimate.ID == id {
			return estimate, nil
		}
	}
	return Estimate{}, fmt.Errorf("estimate with id %d not found", id)
}
