# LOMS (Logistics and Order Management System)

Сервис отвечает за учет заказов и логистику.

## createOrder

Создает новый заказ для пользователя из списка переданных товаров.
Товары при этом нужно зарезервировать на складе.

Request
```json
{
    "user": int64,
    "items": [
        {
            "sku": uint32,
            "count": uint16
        }
    ]
}
```

Response
```json
{
    "orderID": int64
}
```

## listOrder

Показывает информацию по заказу

Request
```json
{
    "orderID": int64
}
```

Response
```json
{
    "status": "string" // (new | awaiting payment | failed | payed | cancelled),
    "user": int64,
    "items": [
        {
            "sku": uint32,
            "count": uint16
        }
    ]
}
```

## orderPayed

Помечает заказ оплаченным. Зарезервированные товары должны перейти в статус купленных.

Request
```json
{
    "orderID": int64
}
```

Response
```json
{}
```

## cancelOrder

Отменяет заказ, снимает резерв со всех товаров в заказе.

Request
```json
{
    "orderID": int64
}
```

Response
```json
{}
```

## stocks

Возвращает количество товаров, которые можно купить с разных складов. Если товар был зарезерванован у кого-то в заказе и ждет оплаты, его купить нельзя.

Request
```json
{
    "sku": uint32
}
```

Response
```json
{
    "stocks": [
        {
            "warehouseID": int64,
            "count": uint64
        }
    ]
}
```

# Checkout

Сервис отвечает за корзину и оформление заказа.

## addToCart

Добавить товар в корзину определенного пользователя. При этом надо проверить наличие товара через `LOMS.stocks`.

Request
```json
{
    "user": int64,
    "sku": uint32,
    "count": uint16
}
```

Response
```json
{}
```

## deleteFromCart

Удалить товар из корзины определенного пользователя.

Request
```json
{
    "user": int64,
    "sku": uint32,
    "count": uint16
}
```

Response
```json
{}
```

## listCart

Показать список товаров в корзине с именами и ценами (их надо в реальном времени получать из ProductService)

Request
```json
{
    "user": int64
}
```

Response
```json
{
    "items": [
        {
            "sku": uint32,
            "count": uint16,
            "name": "string",
            "price": uint32
        }
    ],
    "totalPrice": uint32
}
```

## purchase

Оформить заказ по всем товарам корзины. Вызывает `LOMS.createOrder`.

Request
```json
{
    "user": int64
}
```

Response
```json
{}
```

# Notifications

Будет слушать Кафку и отправлять уведомления, внешнего API нет.

# ProductService

Swagger развернут по адресу:
http://route256.pavl.uk:8080/docs/

GRPC развернуто по адресу:
route256.pavl.uk:8082

## get_product

Request
```json
{
    "token": "string",
    "sku": uint32
}
```

Response
```json
{
    "name": "string",
    "price": uint32
}
```

## list_skus

Request
```json
{
    "token": "string",
    "startAfterSku": uint32,
    "count": uint32
}
```

Response
```json
{
    "skus": [uint32]
}
```

# Путь покупки товаров

- `Checkout.addToCart`
    + добавляем в корзину и проверяем, что есть в наличии
- Можем удалять из корзины
- Можем получать список товаров корзины
    + название и цена тянутся из `ProductService.get_product`
- Приобретаем товары через `Checkout.purchase`
    + идем в `LOMS.createOrder` и создаем заказ
    + У заказа статус `new`
    + `LOMS` резервирует нужное количество единиц товара
    + Если не удалось зарезервить, заказ падает в статус `failed`
    + Если удалось, падаем в статус `awaiting payment`
- Оплачиваем заказ
    + Вызываем `LOMS.orderPayed`
    + Резервы переходят в списание товара со склада
    + Заказ идет в статус `payed`
- Можно отменить заказ до оплаты
    + Вызываем `LOMS.cancelOrder`
    + Все резервирования по заказу отменяются, товары снова доступны другим пользователям
    + Заказ переходит в статус `cancelled`
    + `LOMS` должен сам отменять заказы по таймауту, если не оплатили в течение 10 минут
