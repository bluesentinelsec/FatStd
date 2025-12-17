At the C level, Split/SplitN return an opaque handle fat_StringArray (same underlying handle
  type as fat_String). You can only interact with it through the array functions:

  - Create: fat_StringSplit(s, sep) / fat_StringSplitN(s, sep, n) → fat_StringArray
  - Size: fat_StringArrayLen(a) → size_t
  - Index: fat_StringArrayGet(a, idx) → fat_String (this is a new string handle you must free)
  - Free array: fat_StringArrayFree(a)
  - Join: fat_StringJoin(a, sep) → fat_String (free it)

  Example (C):

  fat_String s = fat_StringNewUTF8("a,b,c");
  fat_String sep = fat_StringNewUTF8(",");
  fat_StringArray a = fat_StringSplit(s, sep);

  size_t n = fat_StringArrayLen(a);
  for (size_t i = 0; i < n; i++) {
    fat_String elem = fat_StringArrayGet(a, i);
    // ... use elem with other fat_String* APIs ...
    fat_StringFree(elem);
  }

  fat_String joined = fat_StringJoin(a, sep);

  fat_StringFree(joined);
  fat_StringArrayFree(a);
  fat_StringFree(sep);
  fat_StringFree(s);

  In Python (ctypes), it’s the same pattern: treat fat_StringArray as a c_size_t handle, call
  fat_StringArrayLen and fat_StringArrayGet, and free each returned fat_String plus the array
  handle when done.

