import IDeliveryLocationVO from "@domains/vos/interfaces/IDeliveryLocationVO";
import IDeliveryProgressVO from "@domains/vos/interfaces/IDeliveryProgressVO";
import IDeliveryStateVO from "@domains/vos/interfaces/IDeliveryStateVO";
import IDeliveryDTO from "@domains/dtos/interfaces/IDeliveryDTO";
export default class DeliveryDTO implements IDeliveryDTO {
    readonly from: IDeliveryLocationVO;
    readonly to: IDeliveryLocationVO;
    readonly progresses: IDeliveryProgressVO[];
    readonly state: IDeliveryStateVO;
    constructor({ from, to, progresses, state }: {
        from: IDeliveryLocationVO;
        to: IDeliveryLocationVO;
        progresses: Array<IDeliveryProgressVO>;
        state: IDeliveryStateVO;
    });
}
//# sourceMappingURL=DeliveryDTO.d.ts.map