package com.lovi.quebic.sockm.context;

import java.util.Collection;
import java.util.List;
import java.util.Map;

public interface ContextListData {

	<E> List<E> getList(String key);

	<E> List<E> getList(String key, Class<E> typeOfValue);
	
	Map<String, List<?>> getDataList();

	void updateListForAdd(String contextKey, Object value);

	void updateListForAdd(String contextKey, int index, Object value);

	void updateListForAddAll(String contextKey, Collection<?> collection);

	void updateListForAddAll(String contextKey, int index,
			Collection<?> collection);

	void updateListForRemove(String contextKey, Object o);

	void updateListForRemoveAll(String contextKey, Collection<?> collection);

	void updateListForClear(String contextKey);

}
