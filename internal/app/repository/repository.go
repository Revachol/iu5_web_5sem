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
	ID    int    // поля структур, которые передаются в шаблон
	Title string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	Price string
	Value string
	Img   string
}

func (r *Repository) GetOrders() ([]Order, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	orders := []Order{ // массив элементов из наших структур
		{
			ID:    1,
			Title: "Кирпич строительный (Российская Империя)",
			Price: "2 руб. 50 коп.",
			Value: "1000 шт.",
			Img:   "/static/img/object_1.jpg",
		},
		{
			ID:    2,
			Title: "Труд каменщика (Викторианская Англия)",
			Price: "2 шиллингов",
			Value: "10 часов",
			Img:   "/static/img/object_2.jpg",
		},
		{
			ID:    3,
			Title: "Мрамор пентелийский (Древняя Греция)",
			Price: "180 драхм",
			Value: "1 талант веса",
			Img:   "/static/img/object_3.jpg",
		},
		{
			ID:    4,
			Title: "Свинец кровельный (Средневековая Франция)",
			Price: "10 денье",
			Value: "лист",
			Img:   "/static/img/object_4.jpg",
		},
		{
			ID:    5,
			Title: "Масло льняное для живописи (Нидерланды, Золотой век)",
			Price: "5 стюверов",
			Value: "унцию",
			Img:   "/static/img/object_5.png",
		},
		{
			ID:    6,
			Title: "Труд разнорабочего (США, Великая Депрессия)",
			Price: "0.50 $",
			Value: "час",
			Img:   "/static/img/object_6.jpg",
		},
	}
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(orders) == 0 {
		return nil, fmt.Errorf("Array is empty")
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
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
