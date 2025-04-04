# delivery-tracker-server

운송장 번호를 통해 택배사를 스크래핑하여 배송 상태를 제공하는 API 서버입니다.  
현재, 웨일 확장 프로그램 "[택배 배송 조회](https://github.com/falsy/delivery-tracker-for-whale)"의 API 서버로 사용되고 있기 때문에 별도로 사용하기 위해서는 몇 가지 수정이 필요합니다.

## Environment

```ts
// .env
DB_USERNAME=
DB_PASSWORD=
DB_NAME=
DB_HOST=
DB_DIALECT=
PORT=
EXTENSION_ID=
DEV_CLIENT_URL=
```

현재 위와 같은 환경 변수를 사용하고 있습니다.  
만약 별도로 사용한다면, CORS 설정을 위한 EXTENSION_ID, DEV_CLIENT_URL를 제거하고 'main.go' 파일에서 CORS 관련 코드를 수정해야 합니다.

## DB

프로젝트에서는 DB로 MySQL을 사용하고 있으며 CarrierModels 라는 테이블에 택배사 정보를 저장하여 사용하고 있습니다.

### 테이블 생성

```ts
CREATE TABLE `CarrierModels` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` char(36) COLLATE utf8mb4_general_ci NOT NULL,
  `no` int NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `displayName` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `isCrawlable` tinyint(1) NOT NULL,
  `isPopupEnabled` tinyint(1) NOT NULL,
  `popupURL` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `createdAt` datetime NOT NULL,
  `updatedAt` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

> • isCrawlable: 서버에서 직접 운송장 추적이 가능한지 여부  
> • isPopupEnabled: 프론트엔드에서 사용자가 새 창으로 운송장 정보를 열 수 있는지 여부  
> • popupURL: 운송장 번호를 붙여 사용할 수 있는 외부 링크

### 데이터 추가

```ts
// 2025-04-05
INSERT INTO `CarrierModels` (`id`, `uid`, `no`, `name`, `displayName`, `isCrawlable`, `isPopupEnabled`, `popupURL`, `createdAt`, `updatedAt`)
VALUES
  (1,'2de90e9c-1fda-11ef-8884-0a8cb08d3aea',1,'epost','우체국 택배',1,1,'https://service.epost.go.kr/trace.RetrieveDomRigiTraceList.comm?displayHeader=N&sid1=','2024-05-31 11:47:09','2024-05-31 11:47:09'),
  (2,'2de91044-1fda-11ef-8884-0a8cb08d3aea',2,'cjlogistics','CJ 대한통운',1,1,'https://trace.cjlogistics.com/next/tracking.html?wblNo=','2024-05-31 11:48:04','2024-05-31 11:48:04'),
  (3,'2de9109b-1fda-11ef-8884-0a8cb08d3aea',3,'hanjin','한진 택배',1,1,'https://www.hanjin.com/kor/CMS/DeliveryMgr/WaybillResult.do?mCode=MN038&schLang=KR&wblnum=','2024-05-31 11:49:03','2024-05-31 11:49:03'),
  (4,'2de910d2-1fda-11ef-8884-0a8cb08d3aea',4,'lotte','롯데 택배',1,1,'https://www.lotteglogis.com/open/tracking?invno=','2024-05-31 11:49:34','2024-05-31 11:49:34'),
  (5,'ed4ec7e2-20ee-11ef-8884-0a8cb08d3aea',5,'logen','로젠 택배',1,1,'https://www.ilogen.com/web/personal/trace/','2024-06-02 14:46:42','2024-06-02 14:46:42'),
  (6,'2de91105-1fda-11ef-8884-0a8cb08d3aea',6,'gspostbox','GS 편의점 택배',0,1,'https://www.cvsnet.co.kr/invoice/tracking.do?invoice_no=','2024-05-31 11:50:45','2024-05-31 11:50:45'),
  (7,'2de9113a-1fda-11ef-8884-0a8cb08d3aea',7,'cupost','CU 편의점 택배',0,1,'https://www.cupost.co.kr/postbox/delivery/localResult.cupost?invoice_no=','2024-05-31 11:51:51','2024-05-31 11:51:51'),
  (8,'2de9116d-1fda-11ef-8884-0a8cb08d3aea',8,'kdexp','경동 택배',1,0,'','2024-05-31 11:52:28','2024-05-31 11:52:28'),
  (9,'2de9119c-1fda-11ef-8884-0a8cb08d3aea',9,'daesin','대신 택배',1,0,'','2024-05-31 11:53:05','2024-05-31 11:53:05'),
  (10,'2de911ca-1fda-11ef-8884-0a8cb08d3aea',10,'ilyanglogis','일량로지스',0,1,'http://www.ilyanglogis.com/functionality/card_form_waybill.asp?hawb_no=','2024-05-31 11:54:03','2024-05-31 11:54:03'),
  (11,'2de911fb-1fda-11ef-8884-0a8cb08d3aea',11,'ems','국체우편(EMS)',0,1,'https://service.epost.go.kr/trace.RetrieveEmsRigiTraceList.comm?displayHeader=N&POST_CODE=','2024-05-31 11:54:36','2024-05-31 11:54:36');
```

## API

### /carriers

서비스에서 제공하는 전체 택배사 정보를 가져옵니다.

```ts
interface ICarrier {
  id: string
  no: number
  name: string
  displayName: string
  isCrawlable: boolean
  isPopupEnabled: boolean
  popupURL: string
}
```

```ts
// e.g.
[
  {
    "id":"2de90e9c-1fda-11ef-8884-0a8cb08d3aea",
    "no":1,
    "name":"epost",
    "displayName":"우체국 택배",
    "isCrawlable":true,
    "isPopupEnabled":true,
    "popupURL":"https://service.epost.go.kr/trace.RetrieveDomRigiTraceList.comm?displayHeader=N&sid1="
  },
  ...
]
```

### /carrier/:carrierId

택배사 id 값을 통해서 해당 택배사 정보를 가져옵니다.

```ts
// e.g.
{
  "id":"2de90e9c-1fda-11ef-8884-0a8cb08d3aea",
  "no":1,
  "name":"epost",
  "displayName":"우체국 택배",
  "isCrawlable":true,
  "isPopupEnabled":true,
  "popupURL":"https://service.epost.go.kr/trace.RetrieveDomRigiTraceList.comm?displayHeader=N&sid1="
}
```

### /delivery/:carrierId/:trackingNumber

택배사 id 값과 운송장 번호를 통해 현재 배송 정보를 가져옵니다.

```ts
interface IDelivery {
  from: IDeliveryLocationVO
  progresses: Array<IDeliveryProgressVO>
  state: IDeliveryStateVO
  to: IDeliveryLocationVO
}

interface IDeliveryLocationVO {
  name: string
  time: string
}

interface IDeliveryProgressVO {
  description: string
  location: string
  time: string
  state: IDeliveryStateVO
}

interface IDeliveryStateVO {
  id: string
  name: string
}
```

```ts
// e.g.
{
  "from": {
    "name": "인천지점",
    "time": "2023-09-12 13:40:00"
  },
  "to": {
    "name": "여의도(대)",
    "time": "2023-09-13 12:00:00"
  },
  "state": {
    "id": "delivered",
    "name": "배달완료"
  },
  "progresses": [
    {
      "description": "배달 완료하였습니다.",
      "location": "여의도(대)",
      "time": "2023-09-13 12:00:00",
      "state": {
        "id": "delivered",
        "name": "배달완료"
      }
    },
    {
      "description": "고객님의 상품을 18~20시에 배달 예정 입니다.",
      "location": "여의도(대)",
      "time": "2023-09-13 09:33:00",
      "state": {
        "id": "out_for_delivery",
        "name": "배달출발"
      }
    },
    ...
  ]
}
```
