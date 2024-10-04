# LOMS (Logistics and Order Management System)

The service is responsible for order management and logistics.

## createOrder

Creates a new order for the user from the list of provided items.
Items need to be reserved in the warehouse.

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

Displays information about an order.

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

Marks the order as paid. Reserved items should change to purchased status.

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

Cancels the order and removes reservations for all items in the order.

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

Returns the quantity of items available for purchase from different warehouses. If an item is reserved in an order awaiting payment, it cannot be purchased.

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

The service is responsible for the shopping cart and order placement.

## addToCart

Adds an item to a specific user's cart. Availability of the item is checked via `LOMS.stocks`.

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

Removes an item from a specific user's cart.

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

Displays a list of items in the cart, including names and prices (retrieved in real time from ProductService).

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

Places an order for all items in the cart by calling `LOMS.createOrder`.

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

Will listen to Kafka and send notifications. No external API.

# ProductService

Swagger available at:
http://route256.pavl.uk:8080/docs/

GRPC available at:
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

# Purchase Flow

- `Checkout.addToCart`
  - Add to cart and check availability
- Items can be removed from the cart
- The list of cart items can be viewed
  - Name and price are retrieved from `ProductService.get_product`
- Purchase items through `Checkout.purchase`
  - Call `LOMS.createOrder` to create an order
  - The order status is `new`
  - `LOMS` reserves the required quantity of items
  - If reservation fails, the order status becomes `failed`
  - If successful, the status becomes `awaiting payment`
- Pay for the order
  - Call `LOMS.orderPayed`
  - Reservations are converted into warehouse deductions
  - The order status changes to `payed`
- The order can be canceled before payment
  - Call `LOMS.cancelOrder`
  - All reservations for the order are canceled, and items are made available to other users again
  - The order status changes to `cancelled`
  - `LOMS` should automatically cancel orders after a timeout if they are not paid within 10 minutes