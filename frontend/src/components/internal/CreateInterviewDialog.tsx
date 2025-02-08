import React, { FunctionComponent } from 'react';
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"

import {MessageCirclePlus} from "lucide-react";
import {Input} from "@/components/ui/input.tsx";
import {Label} from "@/components/ui/label.tsx";

import {TimePicker} from "antd";
import dayjs from "dayjs"

type TProps = {};

const timerFormat = 'HH:mm';

const CreateInterviewDialog: FunctionComponent<TProps> = (props) => {

  return (
      <Dialog>
        <DialogTrigger asChild>
          <Button><MessageCirclePlus /> Новое интервью</Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">

          <DialogHeader>
            <DialogTitle>Создание интервью</DialogTitle>
          </DialogHeader>
          <div>
              <Label className={"pl-1"}>Название интервью</Label>
              <Input className={"mt-1"} name={"title"} placeholder={"Собеседование в Google"}></Input>

              <TimePicker format={timerFormat} />
          </div>
          <DialogFooter>
            <Button type="submit">Save changes</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
  );
};

export default CreateInterviewDialog;
