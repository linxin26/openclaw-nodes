import { useState } from "react"
import {
  Bell,
  CalendarDays,
  Camera,
  Image,
  LocateFixed,
  Monitor,
  PersonStanding,
  Send,
} from "lucide-react"

import { ResultPanel } from "@/components/shared/ResultPanel"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useDeviceOps } from "@/hooks/useDeviceOps"
import { copy } from "@/lib/copy"

export function Operations() {
  const { loading, result, invoke, retry, saveFile } = useDeviceOps()
  const [notificationText, setNotificationText] = useState("OpenClaw GUI 测试通知")
  const [smsTo, setSmsTo] = useState("+15551234567")
  const [smsBody, setSmsBody] = useState("来自 OpenClaw GUI 的测试消息")
  const [calendarTitle, setCalendarTitle] = useState("OpenClaw 评审")
  const [calendarStart, setCalendarStart] = useState("2026-03-29T09:00:00")

  return (
    <div className="space-y-6">
      <Tabs defaultValue="camera">
        <TabsList>
          <TabsTrigger value="camera"><Camera className="mr-2 h-4 w-4" />{copy.pages.operations.tabs.camera}</TabsTrigger>
          <TabsTrigger value="screen"><Monitor className="mr-2 h-4 w-4" />{copy.pages.operations.tabs.screen}</TabsTrigger>
          <TabsTrigger value="location"><LocateFixed className="mr-2 h-4 w-4" />{copy.pages.operations.tabs.location}</TabsTrigger>
          <TabsTrigger value="photos"><Image className="mr-2 h-4 w-4" />{copy.pages.operations.tabs.photos}</TabsTrigger>
          <TabsTrigger value="notifications"><Bell className="mr-2 h-4 w-4" />{copy.pages.operations.tabs.notifications}</TabsTrigger>
          <TabsTrigger value="motion"><PersonStanding className="mr-2 h-4 w-4" />{copy.pages.operations.tabs.motion}</TabsTrigger>
          <TabsTrigger value="sms"><Send className="mr-2 h-4 w-4" />{copy.pages.operations.tabs.sms}</TabsTrigger>
          <TabsTrigger value="calendar"><CalendarDays className="mr-2 h-4 w-4" />{copy.pages.operations.tabs.calendar}</TabsTrigger>
        </TabsList>

        <TabsContent value="camera">
          <CommandCard
            title={copy.pages.operations.camera.title}
            description={copy.pages.operations.camera.description}
            actions={[
              { label: copy.pages.operations.camera.actions.list, run: () => invoke("camera.list") },
              { label: copy.pages.operations.camera.actions.snapshot, run: () => invoke("camera.snap", { cameraId: "0" }) },
              {
                label: copy.pages.operations.camera.actions.clip,
                run: () => invoke("camera.clip", { cameraId: "0", durationMs: 5000 }),
              },
            ]}
            loading={loading}
          />
        </TabsContent>
        <TabsContent value="screen">
          <CommandCard
            title={copy.pages.operations.screen.title}
            description={copy.pages.operations.screen.description}
            actions={[{ label: copy.pages.operations.screen.actions.capture, run: () => invoke("screen.snapshot") }]}
            loading={loading}
          />
        </TabsContent>
        <TabsContent value="location">
          <CommandCard
            title={copy.pages.operations.location.title}
            description={copy.pages.operations.location.description}
            actions={[{ label: copy.pages.operations.location.actions.get, run: () => invoke("location.get") }]}
            loading={loading}
          />
        </TabsContent>
        <TabsContent value="photos">
          <CommandCard
            title={copy.pages.operations.photos.title}
            description={copy.pages.operations.photos.description}
            actions={[{ label: copy.pages.operations.photos.actions.latest, run: () => invoke("photos.latest") }]}
            loading={loading}
          />
        </TabsContent>
        <TabsContent value="notifications">
          <Card>
            <CardHeader>
              <CardTitle>{copy.pages.operations.notifications.title}</CardTitle>
              <CardDescription>{copy.pages.operations.notifications.description}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="notify-body">{copy.pages.operations.notifications.labels.body}</Label>
                <Input
                  id="notify-body"
                  value={notificationText}
                  onChange={(event) => setNotificationText(event.target.value)}
                />
              </div>
              <div className="flex flex-wrap gap-3">
                <Button onClick={() => invoke("notifications.list")} disabled={loading}>
                  {copy.pages.operations.notifications.actions.list}
                </Button>
                <Button
                  variant="secondary"
                  onClick={() => invoke("system.notify", { title: "OpenClaw", body: notificationText })}
                  disabled={loading}
                >
                  {copy.pages.operations.notifications.actions.trigger}
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="motion">
          <CommandCard
            title={copy.pages.operations.motion.title}
            description={copy.pages.operations.motion.description}
            actions={[
              { label: copy.pages.operations.motion.actions.activity, run: () => invoke("motion.activity") },
              { label: copy.pages.operations.motion.actions.pedometer, run: () => invoke("motion.pedometer") },
            ]}
            loading={loading}
          />
        </TabsContent>
        <TabsContent value="sms">
          <Card>
            <CardHeader>
              <CardTitle>{copy.pages.operations.sms.title}</CardTitle>
              <CardDescription>{copy.pages.operations.sms.description}</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="sms-to">{copy.pages.operations.sms.labels.to}</Label>
                <Input id="sms-to" value={smsTo} onChange={(event) => setSmsTo(event.target.value)} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="sms-body">{copy.pages.operations.sms.labels.body}</Label>
                <Input id="sms-body" value={smsBody} onChange={(event) => setSmsBody(event.target.value)} />
              </div>
              <div className="flex flex-wrap gap-3 md:col-span-2">
                <Button onClick={() => invoke("sms.send", { to: smsTo, body: smsBody })} disabled={loading}>
                  {copy.pages.operations.sms.actions.send}
                </Button>
                <Button variant="secondary" onClick={() => invoke("sms.search", { query: smsTo })} disabled={loading}>
                  {copy.pages.operations.sms.actions.search}
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="calendar">
          <Card>
            <CardHeader>
              <CardTitle>{copy.pages.operations.calendar.title}</CardTitle>
              <CardDescription>{copy.pages.operations.calendar.description}</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="calendar-title">{copy.pages.operations.calendar.labels.title}</Label>
                <Input
                  id="calendar-title"
                  value={calendarTitle}
                  onChange={(event) => setCalendarTitle(event.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="calendar-start">{copy.pages.operations.calendar.labels.start}</Label>
                <Input
                  id="calendar-start"
                  value={calendarStart}
                  onChange={(event) => setCalendarStart(event.target.value)}
                />
              </div>
              <div className="flex flex-wrap gap-3 md:col-span-2">
                <Button onClick={() => invoke("calendar.events")} disabled={loading}>
                  {copy.pages.operations.calendar.actions.list}
                </Button>
                <Button
                  variant="secondary"
                  onClick={() => invoke("calendar.add", { title: calendarTitle, start: calendarStart })}
                  disabled={loading}
                >
                  {copy.pages.operations.calendar.actions.add}
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <ResultPanel result={result} onRetry={retry} onSave={saveFile} />
    </div>
  )
}

function CommandCard({
  title,
  description,
  actions,
  loading,
}: {
  title: string
  description: string
  actions: Array<{ label: string; run: () => Promise<unknown> }>
  loading: boolean
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-wrap gap-3">
        {actions.map((action) => (
          <Button key={action.label} onClick={action.run} disabled={loading}>
            {action.label}
          </Button>
        ))}
      </CardContent>
    </Card>
  )
}
