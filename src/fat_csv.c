#include "fat/csv.h"

#include "fatstd_go.h"

fat_CsvReader fat_CsvReaderNewBytes(fat_Bytes data) {
  return (fat_CsvReader)fatstd_go_csv_reader_new_bytes((uintptr_t)data);
}

void fat_CsvReaderFree(fat_CsvReader r) {
  fatstd_go_csv_reader_free((uintptr_t)r);
}

fat_Status fat_CsvReaderRead(fat_CsvReader r, fat_StringArray *out_record, bool *out_eof, fat_Error *out_err) {
  return (fat_Status)fatstd_go_csv_reader_read((uintptr_t)r, (uintptr_t *)out_record, (_Bool *)out_eof, (uintptr_t *)out_err);
}

void fat_CsvReaderFieldPos(fat_CsvReader r, int field, int *out_line, int *out_column) {
  fatstd_go_csv_reader_field_pos((uintptr_t)r, field, out_line, out_column);
}

int64_t fat_CsvReaderInputOffset(fat_CsvReader r) {
  return (int64_t)fatstd_go_csv_reader_input_offset((uintptr_t)r);
}

fat_CsvWriter fat_CsvWriterNewToBytesBuffer(fat_BytesBuffer dst) {
  return (fat_CsvWriter)fatstd_go_csv_writer_new_to_bytes_buffer((uintptr_t)dst);
}

fat_Status fat_CsvWriterWriteRecord(
  fat_CsvWriter w,
  const fat_String *fields,
  size_t n_fields,
  fat_Error *out_err
) {
  return (fat_Status)fatstd_go_csv_writer_write_record((uintptr_t)w, (uintptr_t *)fields, n_fields, (uintptr_t *)out_err);
}

void fat_CsvWriterFlush(fat_CsvWriter w) {
  fatstd_go_csv_writer_flush((uintptr_t)w);
}

fat_Status fat_CsvWriterError(fat_CsvWriter w, fat_Error *out_err) {
  return (fat_Status)fatstd_go_csv_writer_error((uintptr_t)w, (uintptr_t *)out_err);
}

void fat_CsvWriterFree(fat_CsvWriter w) {
  fatstd_go_csv_writer_free((uintptr_t)w);
}
