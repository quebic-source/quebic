package com.lovi.quebic.sockm.context;

import java.util.Map;

public interface ContextMapData {

	<K, V> Map<K, V> getMap(String key);
	
	<K, V> Map<K, V> getMap(String key, Class<K> typeOfKey, Class<V> typeOfValue);

	void updateMapForPut(String contextKey, Object key, Object value);

	void updateMapForPutAll(String contextKey, Map<?, ?> m);

	void updateMapForRemove(String contextKey, Object key);

	void updateMapForClear(String contextKey);

	Map<String, Map<?, ?>> getDataMap();

}
