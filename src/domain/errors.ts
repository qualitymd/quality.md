import * as Schema from "effect/Schema"

const ErrorTypeId = "qualitymd/error"

export class UsageError extends Schema.TaggedErrorClass<UsageError>(`${ErrorTypeId}/UsageError`)(
  "UsageError",
  { detail: Schema.String },
) {
  override get message() {
    return this.detail
  }
}

export class ModelInvalid extends Schema.TaggedErrorClass<ModelInvalid>(
  `${ErrorTypeId}/ModelInvalid`,
)("ModelInvalid", { detail: Schema.String }) {
  override get message() {
    return this.detail
  }
}

export class FileSystemFailure extends Schema.TaggedErrorClass<FileSystemFailure>(
  `${ErrorTypeId}/FileSystemFailure`,
)("FileSystemFailure", { detail: Schema.String }) {
  override get message() {
    return this.detail
  }
}

export class InternalFailure extends Schema.TaggedErrorClass<InternalFailure>(
  `${ErrorTypeId}/InternalFailure`,
)("InternalFailure", { detail: Schema.String }) {
  override get message() {
    return this.detail
  }
}

export type QualitymdError = UsageError | ModelInvalid | FileSystemFailure | InternalFailure
