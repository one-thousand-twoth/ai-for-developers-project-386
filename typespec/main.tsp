import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

@service(#{ title: "Calendar API" })
namespace Calendar;

model Health { status: "ok"; }
model BookingInput {
  @minLength(1) title: string;
  startsAt: utcDateTime;
  endsAt: utcDateTime;
}
model Booking extends BookingInput { id: string; }
@error
model BadRequest {
  @statusCode statusCode: 400;
  @body body: ApiError;
}
@error
model NotFound {
  @statusCode statusCode: 404;
  @body body: ApiError;
}
model ApiError { code: string; message: string; }

@doc("A meeting type. Names are trimmed with Unicode whitespace and compared using Unicode case folding; internal whitespace is significant.")
model MeetingTypeInput {
  @doc("1–200 Unicode code points, with a non-whitespace character. Leading and trailing Unicode whitespace is removed before storing.")
  @minLength(1) @maxLength(200) name: string;
  @minValue(1) @maxValue(480) durationMinutes: int32;
  @doc("Optional description, at most 2000 Unicode code points. Stored unchanged.")
  @maxLength(2000) description?: string;
}
model MeetingType extends MeetingTypeInput { id: string; }

@route("/api/meeting-types") @get @operationId("getMeetingTypes")
op getMeetingTypes(): MeetingType[] | TemporarilyUnavailable;

@route("/api/owner/meeting-types") @post @operationId("createMeetingType")
op createMeetingType(@body meetingType: MeetingTypeInput): {
  @statusCode statusCode: 201;
  @body body: MeetingType;
} | BadRequest | DuplicateMeetingType | TemporarilyUnavailable;

@error
model DuplicateMeetingType {
  @statusCode statusCode: 409;
  @body body: { code: "duplicate_meeting_type"; message: string; };
}
@error
model TemporarilyUnavailable {
  @statusCode statusCode: 503;
  @body body: { code: "temporarily_unavailable"; message: string; };
}

@doc("Deployment readiness endpoint for external container health checks.")
@route("/") @get @operationId("getRoot")
op getRoot(): Health;

@route("/health") @get @operationId("getHealth")
op getHealth(): Health;

@route("/api/bookings") @post @operationId("createBooking")
op createBooking(@body booking: BookingInput): {
  @statusCode statusCode: 201;
  @body body: Booking;
} | BadRequest;

@route("/api/bookings/{id}") @get @operationId("getBooking")
op getBooking(@path id: string): Booking | BadRequest | NotFound;
