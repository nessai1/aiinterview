import React from "react";
import './StatusBadge.css';

enum Status {
    Active = 'active',
    End = 'end',
}

type TProps = {
    status: Status
    children: React.ReactNode
}

const StatusBadge: React.FC<TProps> = (props: TProps) => {
    return <div className={'status-badge ' + props.status}>{props.children}</div>
}

export { StatusBadge, Status };